-- ============================================================================
-- Sistema de Gestión de Laboratorio Clínico
-- Base de Datos: medical_exams_db
-- PostgreSQL 14+
-- ============================================================================

-- Habilitar extensiones necesarias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- TABLA: roles
-- Descripción: Roles y permisos del sistema
-- ============================================================================
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    permissions JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE roles IS 'Roles del sistema con sus permisos';
COMMENT ON COLUMN roles.permissions IS 'Permisos en formato JSON: {"users": ["read", "write"], "orders": ["read"]}';

-- ============================================================================
-- TABLA: users
-- Descripción: Usuarios del sistema (bioanalistas, admin, recepcionistas)
-- ============================================================================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(150) NOT NULL,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    deleted_by INTEGER REFERENCES users(id)
);

COMMENT ON TABLE users IS 'Usuarios del sistema con autenticación';
COMMENT ON COLUMN users.password_hash IS 'Hash bcrypt de la contraseña';
COMMENT ON COLUMN users.locked_until IS 'Bloqueo temporal por intentos fallidos';

-- ============================================================================
-- TABLA: patients
-- Descripción: Información de pacientes
-- ============================================================================
CREATE TABLE patients (
    id SERIAL PRIMARY KEY,
    document_type VARCHAR(20) NOT NULL CHECK (document_type IN ('cedula', 'pasaporte', 'rif', 'otro')),
    document_number VARCHAR(50) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender CHAR(1) CHECK (gender IN ('M', 'F', 'O')),
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100) DEFAULT 'Venezuela',
    emergency_contact_name VARCHAR(150),
    emergency_contact_phone VARCHAR(20),
    blood_type VARCHAR(5) CHECK (blood_type IN ('A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-')),
    allergies TEXT,
    medical_conditions TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER REFERENCES users(id),
    deleted_at TIMESTAMP,
    deleted_by INTEGER REFERENCES users(id),
    UNIQUE(document_type, document_number)
);

COMMENT ON TABLE patients IS 'Registro de pacientes del laboratorio';
COMMENT ON COLUMN patients.allergies IS 'Alergias conocidas del paciente';
COMMENT ON COLUMN patients.medical_conditions IS 'Condiciones médicas relevantes';

-- ============================================================================
-- TABLA: exam_categories
-- Descripción: Categorías de exámenes (Hematología, Serología, etc.)
-- ============================================================================
CREATE TABLE exam_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    code VARCHAR(20) UNIQUE NOT NULL,
    display_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE exam_categories IS 'Categorías de exámenes de laboratorio';

-- ============================================================================
-- TABLA: sample_types
-- Descripción: Tipos de muestras (sangre, orina, heces, etc.)
-- ============================================================================
CREATE TABLE sample_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    collection_instructions TEXT,
    storage_requirements TEXT,
    storage_temperature VARCHAR(50),
    max_storage_time_hours INTEGER,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE sample_types IS 'Tipos de muestras biológicas';
COMMENT ON COLUMN sample_types.collection_instructions IS 'Instrucciones para toma de muestra';
COMMENT ON COLUMN sample_types.storage_requirements IS 'Requisitos de almacenamiento';

-- ============================================================================
-- TABLA: exam_types
-- Descripción: Tipos de exámenes disponibles
-- ============================================================================
CREATE TABLE exam_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    category_id INTEGER NOT NULL REFERENCES exam_categories(id) ON DELETE RESTRICT,
    sample_type_id INTEGER NOT NULL REFERENCES sample_types(id) ON DELETE RESTRICT,
    base_price DECIMAL(10, 2) NOT NULL CHECK (base_price >= 0),
    preparation_instructions TEXT,
    processing_time_hours INTEGER DEFAULT 24,
    requires_fasting BOOLEAN DEFAULT false,
    fasting_hours INTEGER,
    requires_appointment BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE exam_types IS 'Catálogo de tipos de exámenes';
COMMENT ON COLUMN exam_types.code IS 'Código único del examen (ej: VDRL-001)';
COMMENT ON COLUMN exam_types.processing_time_hours IS 'Tiempo estimado de procesamiento';

-- ============================================================================
-- TABLA: exam_parameters
-- Descripción: Parámetros medidos en cada tipo de examen
-- ============================================================================
CREATE TABLE exam_parameters (
    id SERIAL PRIMARY KEY,
    exam_type_id INTEGER NOT NULL REFERENCES exam_types(id) ON DELETE CASCADE,
    parameter_name VARCHAR(200) NOT NULL,
    parameter_code VARCHAR(50),
    unit_of_measure VARCHAR(50),
    reference_min DECIMAL(12, 4),
    reference_max DECIMAL(12, 4),
    reference_value_text TEXT,
    data_type VARCHAR(20) NOT NULL CHECK (data_type IN ('numeric', 'text', 'boolean', 'select')),
    select_options JSONB,
    display_order INTEGER DEFAULT 0,
    is_critical BOOLEAN DEFAULT false,
    is_required BOOLEAN DEFAULT true,
    validation_rules JSONB,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(exam_type_id, parameter_code)
);

COMMENT ON TABLE exam_parameters IS 'Parámetros que se miden en cada examen';
COMMENT ON COLUMN exam_parameters.data_type IS 'Tipo de dato: numeric, text, boolean, select';
COMMENT ON COLUMN exam_parameters.select_options IS 'Opciones para tipo select: ["Reactivo", "No Reactivo"]';
COMMENT ON COLUMN exam_parameters.is_critical IS 'Indica si el parámetro es crítico para alertas';

-- ============================================================================
-- TABLA: reference_ranges
-- Descripción: Rangos de referencia por edad y género
-- ============================================================================
CREATE TABLE reference_ranges (
    id SERIAL PRIMARY KEY,
    exam_parameter_id INTEGER NOT NULL REFERENCES exam_parameters(id) ON DELETE CASCADE,
    gender CHAR(1) CHECK (gender IN ('M', 'F', 'A')) DEFAULT 'A',
    age_min INTEGER DEFAULT 0,
    age_max INTEGER DEFAULT 150,
    reference_min DECIMAL(12, 4),
    reference_max DECIMAL(12, 4),
    reference_text TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE reference_ranges IS 'Rangos de referencia específicos por edad y género';
COMMENT ON COLUMN reference_ranges.gender IS 'M=Masculino, F=Femenino, A=Ambos';

-- ============================================================================
-- TABLA: orders
-- Descripción: Órdenes de exámenes
-- ============================================================================
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    patient_id INTEGER NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'pendiente' 
        CHECK (status IN ('pendiente', 'en_proceso', 'completada', 'cancelada', 'parcial')),
    priority VARCHAR(20) DEFAULT 'normal' 
        CHECK (priority IN ('normal', 'urgente', 'stat')),
    referring_doctor VARCHAR(150),
    doctor_phone VARCHAR(20),
    diagnosis TEXT,
    clinical_notes TEXT,
    subtotal DECIMAL(10, 2) DEFAULT 0,
    discount_percentage DECIMAL(5, 2) DEFAULT 0 CHECK (discount_percentage >= 0 AND discount_percentage <= 100),
    discount_amount DECIMAL(10, 2) DEFAULT 0,
    tax_percentage DECIMAL(5, 2) DEFAULT 0,
    tax_amount DECIMAL(10, 2) DEFAULT 0,
    total_amount DECIMAL(10, 2) DEFAULT 0,
    paid_amount DECIMAL(10, 2) DEFAULT 0,
    balance DECIMAL(10, 2) DEFAULT 0,
    payment_status VARCHAR(20) DEFAULT 'pendiente'
        CHECK (payment_status IN ('pendiente', 'parcial', 'pagado', 'vencido')),
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancelled_by INTEGER REFERENCES users(id),
    cancellation_reason TEXT
);

COMMENT ON TABLE orders IS 'Órdenes de exámenes solicitadas';
COMMENT ON COLUMN orders.order_number IS 'Número de orden generado automáticamente';
COMMENT ON COLUMN orders.priority IS 'Prioridad: normal, urgente, stat (inmediato)';

-- ============================================================================
-- TABLA: order_exams
-- Descripción: Exámenes específicos dentro de cada orden
-- ============================================================================
CREATE TABLE order_exams (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    exam_type_id INTEGER NOT NULL REFERENCES exam_types(id) ON DELETE RESTRICT,
    status VARCHAR(20) NOT NULL DEFAULT 'pendiente'
        CHECK (status IN ('pendiente', 'muestra_tomada', 'en_analisis', 'completado', 'cancelado', 'rechazado')),
    sample_collected_at TIMESTAMP,
    sample_collected_by INTEGER REFERENCES users(id),
    sample_barcode VARCHAR(100),
    analyzed_at TIMESTAMP,
    analyzed_by INTEGER REFERENCES users(id),
    validated_at TIMESTAMP,
    validated_by INTEGER REFERENCES users(id),
    price DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(10, 2) DEFAULT 0,
    final_price DECIMAL(10, 2) NOT NULL,
    notes TEXT,
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE order_exams IS 'Exámenes individuales dentro de cada orden';
COMMENT ON COLUMN order_exams.sample_barcode IS 'Código de barras de la muestra';
COMMENT ON COLUMN order_exams.rejection_reason IS 'Razón de rechazo de muestra';

-- ============================================================================
-- TABLA: exam_results
-- Descripción: Resultados de los exámenes
-- ============================================================================
CREATE TABLE exam_results (
    id SERIAL PRIMARY KEY,
    order_exam_id INTEGER NOT NULL REFERENCES order_exams(id) ON DELETE CASCADE,
    exam_parameter_id INTEGER NOT NULL REFERENCES exam_parameters(id) ON DELETE RESTRICT,
    value_numeric DECIMAL(12, 4),
    value_text TEXT,
    value_boolean BOOLEAN,
    is_abnormal BOOLEAN DEFAULT false,
    is_critical BOOLEAN DEFAULT false,
    abnormality_type VARCHAR(20) CHECK (abnormality_type IN ('low', 'high', 'abnormal', NULL)),
    technician_notes TEXT,
    flags VARCHAR(50),
    version INTEGER DEFAULT 1,
    is_current BOOLEAN DEFAULT true,
    entered_by INTEGER NOT NULL REFERENCES users(id),
    entered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    validated_by INTEGER REFERENCES users(id),
    validated_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(order_exam_id, exam_parameter_id, version)
);

COMMENT ON TABLE exam_results IS 'Resultados de parámetros de exámenes';
COMMENT ON COLUMN exam_results.version IS 'Versión del resultado (para correcciones)';
COMMENT ON COLUMN exam_results.is_current IS 'Indica si es la versión actual';
COMMENT ON COLUMN exam_results.flags IS 'Banderas especiales (L=bajo, H=alto, C=crítico)';

-- ============================================================================
-- TABLA: payments
-- Descripción: Pagos realizados
-- ============================================================================
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    payment_number VARCHAR(50) UNIQUE NOT NULL,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    payment_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    payment_method VARCHAR(30) NOT NULL 
        CHECK (payment_method IN ('efectivo', 'tarjeta_debito', 'tarjeta_credito', 
                                   'transferencia', 'pago_movil', 'cheque', 'otro')),
    reference_number VARCHAR(100),
    bank_name VARCHAR(100),
    card_last_digits VARCHAR(4),
    status VARCHAR(20) DEFAULT 'aprobado'
        CHECK (status IN ('pendiente', 'aprobado', 'rechazado', 'anulado')),
    notes TEXT,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    approved_by INTEGER REFERENCES users(id),
    approved_at TIMESTAMP,
    cancelled_by INTEGER REFERENCES users(id),
    cancelled_at TIMESTAMP,
    cancellation_reason TEXT
);

COMMENT ON TABLE payments IS 'Registro de pagos realizados';
COMMENT ON COLUMN payments.payment_number IS 'Número de recibo generado';

-- ============================================================================
-- TABLA: invoices
-- Descripción: Facturas generadas
-- ============================================================================
CREATE TABLE invoices (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    patient_id INTEGER NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    invoice_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE,
    subtotal DECIMAL(10, 2) NOT NULL,
    discount_amount DECIMAL(10, 2) DEFAULT 0,
    tax_percentage DECIMAL(5, 2) DEFAULT 0,
    tax_amount DECIMAL(10, 2) DEFAULT 0,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pendiente'
        CHECK (status IN ('pendiente', 'pagada', 'parcial', 'vencida', 'anulada')),
    notes TEXT,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    cancelled_by INTEGER REFERENCES users(id),
    cancelled_at TIMESTAMP,
    cancellation_reason TEXT
);

COMMENT ON TABLE invoices IS 'Facturas emitidas';

-- ============================================================================
-- TABLA: audit_logs
-- Descripción: Registro de auditoría de todas las operaciones
-- ============================================================================
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    record_id INTEGER NOT NULL,
    action VARCHAR(10) NOT NULL CHECK (action IN ('INSERT', 'UPDATE', 'DELETE')),
    old_values JSONB,
    new_values JSONB,
    user_id INTEGER REFERENCES users(id),
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE audit_logs IS 'Auditoría completa de cambios en el sistema';

-- ============================================================================
-- TABLA: reagents
-- Descripción: Control de reactivos e insumos
-- ============================================================================
CREATE TABLE reagents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    manufacturer VARCHAR(150),
    supplier VARCHAR(150),
    lot_number VARCHAR(100),
    expiration_date DATE,
    quantity_available DECIMAL(10, 2) NOT NULL DEFAULT 0,
    minimum_stock DECIMAL(10, 2) DEFAULT 0,
    unit_of_measure VARCHAR(50) NOT NULL,
    cost_per_unit DECIMAL(10, 2),
    storage_location VARCHAR(100),
    storage_temperature VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE reagents IS 'Inventario de reactivos e insumos';

-- ============================================================================
-- TABLA: equipment
-- Descripción: Equipos de laboratorio
-- ============================================================================
CREATE TABLE equipment (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    manufacturer VARCHAR(150),
    model VARCHAR(100),
    serial_number VARCHAR(100) UNIQUE,
    purchase_date DATE,
    warranty_expiration DATE,
    last_maintenance_date DATE,
    next_maintenance_date DATE,
    maintenance_frequency_days INTEGER,
    status VARCHAR(30) DEFAULT 'operativo'
        CHECK (status IN ('operativo', 'mantenimiento', 'reparacion', 'fuera_de_servicio', 'retirado')),
    location VARCHAR(100),
    responsible_user_id INTEGER REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE equipment IS 'Equipos y maquinaria del laboratorio';

-- ============================================================================
-- TABLA: equipment_maintenance
-- Descripción: Historial de mantenimiento de equipos
-- ============================================================================
CREATE TABLE equipment_maintenance (
    id SERIAL PRIMARY KEY,
    equipment_id INTEGER NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    maintenance_type VARCHAR(30) NOT NULL 
        CHECK (maintenance_type IN ('preventivo', 'correctivo', 'calibracion', 'verificacion')),
    maintenance_date DATE NOT NULL,
    performed_by VARCHAR(150),
    technician_company VARCHAR(150),
    description TEXT,
    findings TEXT,
    actions_taken TEXT,
    cost DECIMAL(10, 2),
    next_maintenance_date DATE,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE equipment_maintenance IS 'Historial de mantenimiento de equipos';

-- ============================================================================
-- ÍNDICES PARA OPTIMIZACIÓN
-- ============================================================================

-- Usuarios
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role_id);
CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = true;

-- Pacientes
CREATE INDEX idx_patients_document ON patients(document_type, document_number);
CREATE INDEX idx_patients_name ON patients(last_name, first_name);
CREATE INDEX idx_patients_active ON patients(is_active) WHERE is_active = true;
CREATE INDEX idx_patients_created ON patients(created_at DESC);

-- Órdenes
CREATE INDEX idx_orders_patient ON orders(patient_id);
CREATE INDEX idx_orders_date ON orders(order_date DESC);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_number ON orders(order_number);
CREATE INDEX idx_orders_created_by ON orders(created_by);

-- Exámenes de órdenes
CREATE INDEX idx_order_exams_order ON order_exams(order_id);
CREATE INDEX idx_order_exams_type ON order_exams(exam_type_id);
CREATE INDEX idx_order_exams_status ON order_exams(status);
CREATE INDEX idx_order_exams_barcode ON order_exams(sample_barcode);

-- Resultados
CREATE INDEX idx_exam_results_order_exam ON exam_results(order_exam_id);
CREATE INDEX idx_exam_results_parameter ON exam_results(exam_parameter_id);
CREATE INDEX idx_exam_results_current ON exam_results(is_current) WHERE is_current = true;
CREATE INDEX idx_exam_results_abnormal ON exam_results(is_abnormal) WHERE is_abnormal = true;
CREATE INDEX idx_exam_results_critical ON exam_results(is_critical) WHERE is_critical = true;

-- Pagos
CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_date ON payments(payment_date DESC);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_method ON payments(payment_method);

-- Facturas
CREATE INDEX idx_invoices_order ON invoices(order_id);
CREATE INDEX idx_invoices_patient ON invoices(patient_id);
CREATE INDEX idx_invoices_date ON invoices(invoice_date DESC);
CREATE INDEX idx_invoices_status ON invoices(status);

-- Auditoría
CREATE INDEX idx_audit_logs_table_record ON audit_logs(table_name, record_id);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- Tipos de exámenes
CREATE INDEX idx_exam_types_category ON exam_types(category_id);
CREATE INDEX idx_exam_types_sample ON exam_types(sample_type_id);
CREATE INDEX idx_exam_types_active ON exam_types(is_active) WHERE is_active = true;
CREATE INDEX idx_exam_types_code ON exam_types(code);

-- Parámetros de exámenes
CREATE INDEX idx_exam_parameters_type ON exam_parameters(exam_type_id);
CREATE INDEX idx_exam_parameters_order ON exam_parameters(display_order);

-- Reactivos
CREATE INDEX idx_reagents_code ON reagents(code);
CREATE INDEX idx_reagents_expiration ON reagents(expiration_date);
CREATE INDEX idx_reagents_active ON reagents(is_active) WHERE is_active = true;

-- Equipos
CREATE INDEX idx_equipment_code ON equipment(code);
CREATE INDEX idx_equipment_status ON equipment(status);
CREATE INDEX idx_equipment_next_maintenance ON equipment(next_maintenance_date);

-- ============================================================================
-- TRIGGERS PARA UPDATED_AT AUTOMÁTICO
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Aplicar trigger a todas las tablas con updated_at
CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_patients_updated_at BEFORE UPDATE ON patients
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_exam_categories_updated_at BEFORE UPDATE ON exam_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sample_types_updated_at BEFORE UPDATE ON sample_types
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_exam_types_updated_at BEFORE UPDATE ON exam_types
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_exam_parameters_updated_at BEFORE UPDATE ON exam_parameters
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_order_exams_updated_at BEFORE UPDATE ON order_exams
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_exam_results_updated_at BEFORE UPDATE ON exam_results
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reagents_updated_at BEFORE UPDATE ON reagents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_equipment_updated_at BEFORE UPDATE ON equipment
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- FUNCIÓN PARA GENERAR NÚMERO DE ORDEN
-- ============================================================================

CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.order_number IS NULL THEN
        NEW.order_number := 'ORD-' || TO_CHAR(CURRENT_DATE, 'YYYYMMDD') || '-' || 
                           LPAD(NEXTVAL('orders_id_seq')::TEXT, 6, '0');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER generate_order_number_trigger
BEFORE INSERT ON orders
FOR EACH ROW EXECUTE FUNCTION generate_order_number();

-- ============================================================================
-- FUNCIÓN PARA GENERAR NÚMERO DE PAGO
-- ============================================================================

CREATE OR REPLACE FUNCTION generate_payment_number()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.payment_number IS NULL THEN
        NEW.payment_number := 'PAY-' || TO_CHAR(CURRENT_DATE, 'YYYYMMDD') || '-' || 
                             LPAD(NEXTVAL('payments_id_seq')::TEXT, 6, '0');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER generate_payment_number_trigger
BEFORE INSERT ON payments
FOR EACH ROW EXECUTE FUNCTION generate_payment_number();

-- ============================================================================
-- FUNCIÓN PARA GENERAR NÚMERO DE FACTURA
-- ============================================================================

CREATE OR REPLACE FUNCTION generate_invoice_number()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.invoice_number IS NULL THEN
        NEW.invoice_number := 'INV-' || TO_CHAR(CURRENT_DATE, 'YYYYMMDD') || '-' || 
                             LPAD(NEXTVAL('invoices_id_seq')::TEXT, 6, '0');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER generate_invoice_number_trigger
BEFORE INSERT ON invoices
FOR EACH ROW EXECUTE FUNCTION generate_invoice_number();

-- ============================================================================
-- FUNCIÓN PARA CALCULAR TOTALES DE ORDEN
-- ============================================================================

CREATE OR REPLACE FUNCTION calculate_order_totals()
RETURNS TRIGGER AS $$
DECLARE
    order_subtotal DECIMAL(10,2);
BEGIN
    -- Calcular subtotal de todos los exámenes
    SELECT COALESCE(SUM(final_price), 0) INTO order_subtotal
    FROM order_exams
    WHERE order_id = NEW.order_id AND status != 'cancelado';
    
    -- Actualizar la orden
    UPDATE orders SET
        subtotal = order_subtotal,
        discount_amount = (order_subtotal * discount_percentage / 100),
        tax_amount = ((order_subtotal - (order_subtotal * discount_percentage / 100)) * tax_percentage / 100),
        total_amount = (order_subtotal - (order_subtotal * discount_percentage / 100)) + 
                      ((order_subtotal - (order_subtotal * discount_percentage / 100)) * tax_percentage / 100),
        balance = total_amount - paid_amount
    WHERE id = NEW.order_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER calculate_order_totals_trigger
AFTER INSERT OR UPDATE ON order_exams
FOR EACH ROW EXECUTE FUNCTION calculate_order_totals();

-- ============================================================================
-- FUNCIÓN PARA ACTUALIZAR PAID_AMOUNT EN ORDEN
-- ============================================================================

CREATE OR REPLACE FUNCTION update_order_paid_amount()
RETURNS TRIGGER AS $$
DECLARE
    total_paid DECIMAL(10,2);
BEGIN
    -- Calcular total pagado
    SELECT COALESCE(SUM(amount), 0) INTO total_paid
    FROM payments
    WHERE order_id = NEW.order_id AND status = 'aprobado';
    
    -- Actualizar la orden
    UPDATE orders SET
        paid_amount = total_paid,
        balance = total_amount - total_paid,
        payment_status = CASE
            WHEN total_paid = 0 THEN 'pendiente'
            WHEN total_paid < total_amount THEN 'parcial'
            ELSE 'pagado'
        END
    WHERE id = NEW.order_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_order_paid_amount_trigger
AFTER INSERT OR UPDATE ON payments
FOR EACH ROW EXECUTE FUNCTION update_order_paid_amount();

-- ============================================================================
-- DATOS INICIALES
-- ============================================================================

-- Roles del sistema
INSERT INTO roles (name, description, permissions) VALUES
('admin', 'Administrador del sistema', '{"all": ["*"]}'),
('bioanalista', 'Bioanalista - Análisis y validación', '{"orders": ["read"], "results": ["read", "write"], "patients": ["read"]}'),
('recepcionista', 'Recepcionista - Gestión de pacientes y órdenes', '{"patients": ["read", "write"], "orders": ["read", "write"], "payments": ["read", "write"]}'),
('contador', 'Contador - Gestión financiera', '{"payments": ["read"], "invoices": ["read", "write"], "reports": ["read"]}');

-- Usuario administrador inicial (password: Admin123!)
-- Hash generado con bcrypt cost 10
INSERT INTO users (username, email, password_hash, full_name, role_id) VALUES
('admin', 'admin@laboratorio.com', '$2a$10$YourHashHere', 'Administrador del Sistema', 1);

-- Categorías de exámenes
INSERT INTO exam_categories (name, code, description, display_order) VALUES
('Hematología', 'HEM', 'Estudios de sangre y componentes', 1),
('Química Sanguínea', 'QS', 'Análisis bioquímicos en sangre', 2),
('Serología', 'SER', 'Detección de anticuerpos y antígenos', 3),
('Microbiología', 'MICRO', 'Cultivos y antibiogramas', 4),
('Urianálisis', 'URI', 'Análisis de orina', 5),
('Coprología', 'COPRO', 'Análisis de heces', 6),
('Inmunología', 'INMUNO', 'Pruebas inmunológicas', 7),
('Hormonas', 'HORM', 'Perfil hormonal', 8),
('Marcadores Tumorales', 'MT', 'Detección de marcadores de cáncer', 9);

-- Tipos de muestras
INSERT INTO sample_types (name, description, collection_instructions, storage_temperature, max_storage_time_hours) VALUES
('Sangre Venosa', 'Muestra de sangre obtenida por venopunción', 'Extraer en tubo apropiado según el análisis', '2-8°C', 24),
('Sangre Capilar', 'Muestra de sangre por punción digital', 'Limpiar área, punzar y recolectar', '2-8°C', 12),
('Orina', 'Muestra de orina (primera de la mañana preferible)', 'Recolectar en frasco estéril', '2-8°C', 24),
('Heces', 'Muestra de materia fecal', 'Recolectar en frasco estéril sin contaminación de orina', '2-8°C', 48),
('Hisopado Nasal', 'Muestra de mucosa nasal', 'Introducir hisopo en ambas fosas nasales', '2-8°C', 48),
('Hisopado Faríngeo', 'Muestra de mucosa faríngea', 'Frotar ambas amígdalas y pared posterior', '2-8°C', 48),
('Esputo', 'Muestra de secreción bronquial', 'Toser profundamente y expectorar', '2-8°C', 24);

-- ============================================================================
-- VISTAS ÚTILES
-- ============================================================================

-- Vista de órdenes con información completa
CREATE OR REPLACE VIEW v_orders_complete AS
SELECT 
    o.id,
    o.order_number,
    o.order_date,
    o.status,
    o.priority,
    p.id as patient_id,
    p.document_number,
    p.first_name || ' ' || p.last_name as patient_name,
    p.phone as patient_phone,
    u.full_name as created_by_name,
    o.total_amount,
    o.paid_amount,
    o.balance,
    o.payment_status,
    COUNT(oe.id) as total_exams,
    COUNT(CASE WHEN oe.status = 'completado' THEN 1 END) as completed_exams
FROM orders o
JOIN patients p ON o.patient_id = p.id
JOIN users u ON o.created_by = u.id
LEFT JOIN order_exams oe ON o.id = oe.order_id
GROUP BY o.id, p.id, u.id;

-- Vista de resultados pendientes de validación
CREATE OR REPLACE VIEW v_pending_validations AS
SELECT 
    oe.id as order_exam_id,
    o.order_number,
    p.first_name || ' ' || p.last_name as patient_name,
    et.name as exam_name,
    oe.analyzed_at,
    u.full_name as analyzed_by,
    COUNT(er.id) as total_results,
    COUNT(CASE WHEN er.validated_at IS NULL THEN 1 END) as pending_results
FROM order_exams oe
JOIN orders o ON oe.order_id = o.id
JOIN patients p ON o.patient_id = p.id
JOIN exam_types et ON oe.exam_type_id = et.id
JOIN users u ON oe.analyzed_by = u.id
LEFT JOIN exam_results er ON oe.id = er.order_exam_id AND er.is_current = true
WHERE oe.status = 'en_analisis'
GROUP BY oe.id, o.order_number, p.first_name, p.last_name, et.name, oe.analyzed_at, u.full_name
HAVING COUNT(CASE WHEN er.validated_at IS NULL THEN 1 END) > 0;

COMMENT ON VIEW v_orders_complete IS 'Vista completa de órdenes con información del paciente y estadísticas';
COMMENT ON VIEW v_pending_validations IS 'Vista de resultados pendientes de validación';

-- ============================================================================
-- PERMISOS FINALES
-- ============================================================================

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO lab_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO lab_admin;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO lab_admin;

-- ============================================================================
-- FIN DEL SCRIPT
-- ============================================================================