CREATE TABLE Usuario (
    id_usuario SERIAL PRIMARY KEY,              -- Identificador interno
    nombre VARCHAR(100) NOT NULL,               -- Nombre del usuario
    apellido VARCHAR(100) NOT NULL,             -- Apellido del usuario
    email VARCHAR(100) UNIQUE NOT NULL,         -- Correo electrónico único 
    contraseña VARCHAR(255) NOT NULL,           -- Contraseña (hashed)
    fecha_registro TIMESTAMP WITH TIME ZONE 
        DEFAULT CURRENT_TIMESTAMP               -- Fecha de registro
);

CREATE TABLE Movimiento (
    id_movimiento SERIAL PRIMARY KEY,           -- Identificador único
    id_usuario INT NOT NULL,                    -- Relación con Usuario
    monto FLOAT NOT NULL CHECK (monto > 0),               -- Monto 
    tipo VARCHAR(1) NOT NULL CHECK (tipo IN ('I','G')), -- (I) ingreso, (G) gasto
    descripcion VARCHAR(400),                           -- Descripción opcional
    fecha_movimiento DATE NOT NULL,             -- Fecha del movimiento (cargada por el usuario)
    FOREIGN KEY (id_usuario) REFERENCES Usuario(id_usuario)
);

INSERT INTO Usuario (nombre, apellido, email, contraseña)
VALUES
('Ana', 'García', 'ana.garcia@example.com', 'hash123'),
('Bruno', 'López', 'bruno.lopez@example.com', 'hash456'),
('Carla', 'Pérez', 'carla.perez@example.com', 'hash789');
INSERT INTO Movimiento (id_usuario, monto, tipo, descripcion, fecha_movimiento)
VALUES
-- Movimientos de Ana
(1, 1200.00, 'I', 'Sueldo mensual', '2025-10-01'),
(1, 300.00, 'G', 'Compra supermercado', '2025-10-03'),
(1, 150.00, 'G', 'Cine con amigos', '2025-10-05'),

-- Movimientos de Bruno
(2, 2500.00, 'I', 'Pago freelance', '2025-10-02'),
(2, 500.00, 'G', 'Cena en restaurante', '2025-10-04'),
(2, 800.00, 'G', 'Reparación de computadora', '2025-10-06'),

-- Movimientos de Carla
(3, 1800.00, 'I', 'Venta de productos', '2025-10-01'),
(3, 400.00, 'G', 'Ropa nueva', '2025-10-03'),
(3, 250.00, 'G', 'Transporte', '2025-10-05');
