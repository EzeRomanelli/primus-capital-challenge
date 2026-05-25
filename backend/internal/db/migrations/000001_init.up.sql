-- Schema inicial de Northwind Cobranza.
-- 4 tablas. Postgres 16 trae gen_random_uuid() en core.

-- ============================================================================
-- segmentos: tabla de configuracion con 4 filas fijas.
-- tolerancia_dias = dias post-vencimiento antes de pesar en el scoring.
-- ============================================================================
CREATE TABLE segmentos (
    nombre TEXT PRIMARY KEY,
    tolerancia_dias INT NOT NULL CHECK (tolerancia_dias >= 0),
    descripcion TEXT NOT NULL
);

INSERT INTO segmentos (nombre, tolerancia_dias, descripcion) VALUES
    ('corporativo', 30, 'MRR alto y payment terms acordados largos. Tolera 30 dias post-vencimiento.'),
    ('pyme_sana',   15, 'Historial puntual de pago. Tolera 15 dias post-vencimiento.'),
    ('en_riesgo',    5, 'Empezo a atrasarse vs su patron historico. Tolera 5 dias post-vencimiento.'),
    ('zombi',        0, '90+ dias sin pagar y sigue consumiendo. Cero tolerancia.');

-- ============================================================================
-- clientes
-- segmento = calculado por el suggester (reglas explicitas). Sin override
-- manual en el MVP: si la analista necesita reclasificar, lo hacemos por
-- iteracion 2.
-- ============================================================================
CREATE TABLE clientes (
    id                  UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre              TEXT         NOT NULL,
    industria           TEXT,
    fecha_alta          DATE         NOT NULL,
    mrr_usd             NUMERIC(12,2) NOT NULL CHECK (mrr_usd >= 0),
    payment_terms_dias  INT          NOT NULL CHECK (payment_terms_dias > 0),
    segmento            TEXT         NOT NULL REFERENCES segmentos(nombre),
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_clientes_segmento ON clientes(segmento);

-- ============================================================================
-- facturas
-- estado puede ser 'pendiente' | 'pagada' | 'vencida'.
-- ============================================================================
CREATE TABLE facturas (
    id                  UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    cliente_id          UUID          NOT NULL REFERENCES clientes(id),
    numero              TEXT          NOT NULL,
    fecha_emision       DATE          NOT NULL,
    fecha_vencimiento   DATE          NOT NULL,
    fecha_pago          DATE,
    monto_usd           NUMERIC(12,2) NOT NULL CHECK (monto_usd >= 0),
    estado              TEXT          NOT NULL CHECK (estado IN ('pendiente','pagada','vencida')),
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX idx_facturas_cliente_id ON facturas(cliente_id);
CREATE INDEX idx_facturas_estado     ON facturas(estado);

-- ============================================================================
-- gestiones: registro de cada interaccion de cobranza.
-- ============================================================================
CREATE TABLE gestiones (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    cliente_id  UUID         NOT NULL REFERENCES clientes(id),
    fecha       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    tipo        TEXT         NOT NULL CHECK (tipo IN ('llamada','email','whatsapp','visita')),
    resultado   TEXT         NOT NULL CHECK (resultado IN ('sin_respuesta','promesa_pago','disputa','pagado','otro')),
    notas       TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_gestiones_cliente_fecha ON gestiones(cliente_id, fecha DESC);
