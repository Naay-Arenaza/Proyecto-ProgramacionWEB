package handlers

import (
	db "ProyectoFinanzas/db/sqlc"
	"ProyectoFinanzas/logic"
	"ProyectoFinanzas/views"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

type MovimientoWebHandler struct {
	movLogic  *logic.MovCapaLogica
	userLogic *logic.UserCapaLogica
}

func NewMovimientoWebHandler(movL *logic.MovCapaLogica, userL *logic.UserCapaLogica) *MovimientoWebHandler {
	return &MovimientoWebHandler{
		movLogic:  movL,
		userLogic: userL,
	}
}

type contextKey string

const UserIDKey contextKey = "userID"
const SessionCookieName = "auth_session"

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *MovimientoWebHandler) Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		errMsg := r.URL.Query().Get("error")

		if err := views.SigninPage(errMsg).Render(r.Context(), w); err != nil {
			http.Error(w, "Error al renderizar Signin", http.StatusInternalServerError)
		}
		return
	}

	// POST: Lógica de Autenticación
	if err := r.ParseForm(); err != nil {
		log.Printf("Error al parsear formulario de Signin: %v", err)
		http.Redirect(w, r, "/signin?error=Error al procesar la solicitud", http.StatusSeeOther)
		return
	}

	email := r.FormValue("usuario")
	password := r.FormValue("password")

	user, err := h.userLogic.Authenticate(r.Context(), email, password)

	if err != nil {
		log.Printf("Fallo de autenticación para %s: %v", email, err)

		w.Header().Set("HX-Redirect", "/signin?error=Credenciales+inválidas")
		w.WriteHeader(http.StatusOK)
		return
	}

	sessionToken, expiresAt, err := h.userLogic.CreateUserSession(r.Context(), user.IDUsuario)
	if err != nil {
		log.Printf("Error al crear sesión persistente para usuario %d: %v", user.IDUsuario, err)
		http.Error(w, "Error interno al iniciar sesión.", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https",
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (h *MovimientoWebHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc { // Verifica que la sesion este activa cada vez que se hace un cambio en la URL
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(SessionCookieName)
		if err != nil { // No existe la cookie (usuario no autenticado)
			http.Redirect(w, r, "/signin", http.StatusFound) // Si la cookie no existe, redirigimos al login
			return
		}

		sessionToken := c.Value

		userSession, err := h.userLogic.GetUserSession(r.Context(), sessionToken) // Se valida del token y que la sesion siga activa

		if err != nil {
			h.clearAuthCookie(w)
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userSession.IDUsuario) // Guarda el IDUsuario en el contexto de la peticion

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (h *MovimientoWebHandler) clearAuthCookie(w http.ResponseWriter) { //Se encarga de "cerrar la sesión" del lado del cliente invalidando la cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour), // Una fecha en el pasado -> Para que el navegador elimine la cookie inmediatamente.
		HttpOnly: true,                       // Mitiga ataques XSS
		Secure:   true,                       // Solo se envíe sobre conexiones encriptadas (HTTPS)
		SameSite: http.SameSiteLaxMode,       // Mitigar ataques CSRF
	})
}

func (h *MovimientoWebHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		http.Error(w, "No hay sesión activa para refrescar", http.StatusUnauthorized)
		return
	}

	sessionToken := c.Value

	// 1. Validar sesión existente
	userSession, err := h.userLogic.GetUserSession(r.Context(), sessionToken)
	if err != nil {
		h.clearAuthCookie(w)
		http.Error(w, "Sesión inválida o expirada", http.StatusUnauthorized)
		return
	}

	// 2. Borrar la sesión antigua
	h.userLogic.DeleteUserSession(r.Context(), sessionToken)

	// 3. Crear una nueva sesión persistente con nueva expiración
	newSessionToken, expiresAt, err := h.userLogic.CreateUserSession(r.Context(), userSession.IDUsuario)
	if err != nil {
		http.Error(w, "Error interno al crear nueva sesión", http.StatusInternalServerError)
		return
	}

	// 4. Enviar la nueva cookie al cliente
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    newSessionToken,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sesión refrescada con éxito."))
}

func (h *MovimientoWebHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) { // Cerrar Sesion
	c, err := r.Cookie(SessionCookieName)
	if err == nil {
		h.userLogic.DeleteUserSession(r.Context(), c.Value)
	}

	h.clearAuthCookie(w)
	w.Header().Set("HX-Redirect", "/signin")
	w.WriteHeader(http.StatusOK) // Redirección completa a /signin
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *MovimientoWebHandler) EditMovimientoHandler(w http.ResponseWriter, r *http.Request) { // Toma el movimiento a editar pedido por el Usuario
	idStr := strings.TrimPrefix(r.URL.Path, "/movimientos/edit/")
	movimientoID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	userID, err1 := ctx.Value(UserIDKey).(int32) // Sacamos el IDUsario del contexto de la peticion que le mandamos por el Middleware
	if !err1 {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	var MovArg db.GetMovimientoParams

	MovArg.IDMovimiento = int32(movimientoID) // IDMovimiento de la URL
	MovArg.IDUsuario = int32(userID)          // IDUsuario del contexto de la peticion

	mov, err2 := h.movLogic.GetMovimientoLogic(r.Context(), MovArg)
	if err2 != nil {
		if errors.Is(err2, sql.ErrNoRows) {
			http.NotFound(w, r) // Si no se encuentra en la BD (ya sea xq no es tu movimiento o xq no existe), devolvemos 404
			return
		}
		http.Error(w, "Error interno al cargar datos", http.StatusInternalServerError)
		return
	}

	html := views.MovimientoEditForm(mov)
	log.Printf("Cargando Movimiento ID %d, Tipo: '%s'", mov.IDMovimiento, mov.Tipo)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html.Render(r.Context(), w)
}

func (h *MovimientoWebHandler) ServeForm(w http.ResponseWriter, r *http.Request) { // Pagina Inicial -> Formulario de CrearMovimiento y ListMovimientos del Usuario
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	ctx := r.Context()

	userID, err := ctx.Value(UserIDKey).(int32) // Sacamos el IDUsario del contexto de la peticion que le mandamos por el Middleware
	if !err {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	movimientos, err2 := h.movLogic.ListMovimientoLogic(ctx, userID) // Buscamos todos los movimientos del Usuario

	if err2 != nil {
		log.Printf("Error al cargar movimientos: %v", err2)
		http.Error(w, "Error interno del servidor al cargar datos", http.StatusInternalServerError)
		return
	}

	comp := views.Container(movimientos)
	templ.Handler(views.Layout("MovFinanzas", comp)).ServeHTTP(w, r)
}

// /////////////////////////////////////////////// --->  /MOVIMIENTOS
func (q *MovimientoWebHandler) MovimientosHandler(w http.ResponseWriter, r *http.Request) { // POST (CrearMovimiento)
	switch r.Method {
	case http.MethodPost:
		q.PostMovimiento(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ///////////////////////// -> CrearMovimiento
func (h *MovimientoWebHandler) PostMovimiento(w http.ResponseWriter, r *http.Request) {
	var newMovimiento db.CreateMovimientoParams
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parseando formulario", http.StatusBadRequest)
		return
	}

	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	sessionToken := c.Value

	userSession, err1 := h.userLogic.GetUserSession(r.Context(), sessionToken)
	if err1 != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	newMovimiento.IDUsuario = userSession.IDUsuario
	newMovimiento.Tipo = r.FormValue("tipo")
	monto, err1 := strconv.ParseFloat(r.FormValue("monto"), 64)

	if err1 != nil {
		http.Error(w, "Error de conversion del monto", http.StatusBadRequest)
		return
	}
	newMovimiento.Monto = monto

	newMovimiento.Descripcion = sql.NullString{
		String: r.FormValue("descripcion"),
		Valid:  r.FormValue("descripcion") != "",
	}

	fechaStr := r.FormValue("fechaMovimiento")
	fecha, _ := time.Parse("2006-01-02", fechaStr)
	if !logic.EsFechaValida(fecha) {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}
	newMovimiento.FechaMovimiento = fecha

	_, err2 := h.movLogic.CreateMovimientoLogic(r.Context(), newMovimiento)

	if err2 != nil {
		log.Printf("Error de lógica al crear: %v", err2)
		if strings.Contains(err2.Error(), "el monto del movmiento no puede ser menor o igual a 0") ||
			strings.Contains(err2.Error(), "la fecha debe ser menor a la actual") {
			http.Error(w, "Error de validación: "+err2.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Error interno del servidor al guardar: "+err2.Error(), http.StatusInternalServerError)
		return
	}

	movimientos, err3 := h.movLogic.ListMovimientoLogic(r.Context(), userSession.IDUsuario)

	if err3 != nil {
		http.Error(w, "Error al obtener movimientos: "+err3.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err4 := views.MovimientoList(movimientos).Render(r.Context(), w); err4 != nil {
		http.Error(w, "Error al renderizar la lista: "+err4.Error(), http.StatusInternalServerError)
		return
	}
}

// /////////////////////////////////////////////// --->  /MOVIMIENTO/
func (h *MovimientoWebHandler) MovimientoHandler(w http.ResponseWriter, r *http.Request) { // POST (UpdateMovimiento) y Delete (DeleteMov)
	log.Printf("ID recibido: ")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	log.Printf("ID recibido: %d", id)
	if err != nil {
		http.Error(w, "Id de movimiento invalido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	// case http.MethodGet:
	//     h.getMov(w, r, id)
	case http.MethodPost:
		h.updateMovimiento(w, r, id)
	case http.MethodDelete:
		h.deleteMov(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ///////////////////////// -> UPDATE
func (h *MovimientoWebHandler) updateMovimiento(w http.ResponseWriter, r *http.Request, id int) {
	var newMovimiento db.UpdateMovimientoParams
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID, err1 := ctx.Value(UserIDKey).(int32) // Sacamos el IDUsario del contexto de la peticion que le mandamos por el Middleware
	if !err1 {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	newMovimiento.IDUsuario = int32(userID)
	newMovimiento.IDMovimiento = int32(id)
	newMovimiento.Tipo = r.FormValue("tipo")
	montoStr := r.FormValue("monto")
	sanitizedStr := strings.ReplaceAll(montoStr, ",", ".")
	monto, err := strconv.ParseFloat(sanitizedStr, 64)
	if err != nil {
		fmt.Printf("Error de formato: %v. El valor original era: %s\n", err, monto)
		return
	}
	if !logic.MontoValido(monto) {
		http.Error(w, "Monto invalido", http.StatusBadRequest)
		return
	}
	newMovimiento.Monto = monto
	newMovimiento.Descripcion = sql.NullString{
		String: r.FormValue("descripcion"),
		Valid:  r.FormValue("descripcion") != "",
	}
	fechaStr := r.FormValue("fechaMovimiento")
	fecha, _ := time.Parse("2006-01-02", fechaStr)
	if !logic.EsFechaValida(fecha) {
		http.Error(w, "Fecha inválida", http.StatusBadRequest)
		return
	}
	newMovimiento.FechaMovimiento = fecha

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	h.movLogic.UpdateMovimientoLogic(ctx, newMovimiento)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ///////////////////////// -> DELETE
func (h *MovimientoWebHandler) deleteMov(w http.ResponseWriter, r *http.Request, id int) {
	err := h.movLogic.DeleteMovimientoLogic(r.Context(), int32(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
}
