# DECISIONS — Northwind Cobranza MVP

Las decisiones más relevantes que tomé y por qué. Decisor único: yo.

---

## 1. Foco en la urgencia, no anticipación

**Contexto.** La CEO pidió 3 cosas en una frase: gestionar la cobranza, anticiparnos a problemas, saber dónde poner foco.

**Decidí.** El alma del MVP es **priorización + segmentación** ("foco"). Cortar anticipación y dejar gestión en lo mínimo necesario.

**Por qué.** "Foco" es la palabra textual de la CEO. Anticipación creíble requiere histórico real + ML defendible; en 3 días sería un mock disfrazado. La gestión es CRUD: implementarla entera come tiempo sin diferenciar el producto.

**Qué descarté.** Predicción de mora, gráficos de tendencia, vista CEO con KPIs. Todo iteración 2.

---

## 2. Score con 2 factores, no 4

**Decidí.** `score = urgencia × 0.6 + impacto × 0.4`. Urgencia pesa más porque el dolor del MVP es "apagar incendios" — el tiempo aprieta más que el monto.

**Por qué.** Dos factores se explican en una frase y se defienden con números reales (días de atraso, USD pendientes). Cuatro factores requieren párrafo y nadie los recuerda.

**Qué descarté.** Mi primer intento (sugerido por Claude) tenía 4 factores ponderados (urgencia + impacto + descuido + riesgo). Lo corté después de articular el principio de simplicidad.

**Dónde se ve.** `backend/internal/scoring/scoring.go` — función pura de 50 líneas. Tabla de casos en `scoring_test.go`.

---

## 3. Stack mínimo: Go + Chi + pgx + SQL plano (sin ORM)

**Decidí.** Driver puro `pgx`, queries inline en archivos `.go`. Router `chi`. Sin ORM, sin codegen.

**Por qué.** 4 tablas no justifican una capa de abstracción. `pool.Query(ctx, "SELECT ...")` se entiende sin aprender un DSL. Tests más fáciles porque las funciones son `(ctx, pool, args) → result`.

**Qué descarté.** sqlc (codegen — overkill para 4 tablas), GORM (magia con reflection, no idiomático en Go).

---

## 4. Sin override manual del segmento en el MVP

**Decidí.** El segmento se calcula por reglas explícitas (suggester) y se actualiza con un job fire-and-forget cuando hay gestión nueva. **No hay override de la analista** sobre la regla automática.

**Por qué.** El enunciado pide "**un flujo principal usable end-to-end**". Entregar dos flujos a medias (registrar gestión + sobreescribir segmento) era peor que entregar uno sólido (registrar gestión). El override queda como iteración 2 si la analista lo pide después de usar el sistema unas semanas.

**Qué descarté.** Campos `segmento_actual` vs `segmento_sugerido`, `override_motivo`, `override_fecha`. Endpoint PATCH. Dialog y banner en frontend.

**Dónde se ve.** Schema: una sola columna `segmento` en `clientes`. Backend: 5 endpoints (no 6).

---

## 5. Docker compose full-stack: setup de 1 comando

**Decidí.** Todo en containers (Postgres + migrate + seed + backend + frontend) orquestado por `docker compose up`. El único pre-requisito en el host del evaluador es Docker.

**Por qué.** El enunciado pide *"el sistema levanta siguiendo el README en un equipo limpio en menos de 10 minutos"*. Versión sin docker eran 6 pasos + 4 dependencias (Go, Node, npm, golang-migrate). Versión con docker es 2 pasos y cero deps no-Docker.

**Qué descarté.** Setup local "puro" donde el evaluador instala Go/Node. Lo mantengo solo en la sección "Modo desarrollo" del README para mí (iteración con HMR).

---

## 6. No contactar a jrain

**Contexto.** El enunciado autoriza preguntas a `jrain@primuscapital.cl` y aclara que las respuestas llegan en 4-6h. Recibí el documento un viernes.

**Decidí.** No mandar ninguna pregunta. Asumí 10 cosas razonables y las dejé explícitas.

**Por qué.** El enunciado dice textual *"confío en tu criterio"* y *"si algo te parece raro lo dejamos así a propósito"* — es invitación a decidir, no a preguntar. Mis supuestos no afectan el alma del producto. En un proyecto real con cliente real hubiese preguntado 2-3 cosas críticas el viernes; en un challenge de 3 días con un evaluador que pide criterio bajo ambigüedad, las preguntas defensivas son la opción inferior.

**Supuestos asumidos:** single-user (sin auth), moneda única USD, sin envío de mails reales, sin importador de Sheet (seed sintético), single-tenant, payment_terms fijos por cliente, fechas UTC en backend con locale `es-AR` en frontend, sin paginación (420 caben), recálculo del segmento por job fire-and-forget al crear gestión.

---

## 7. Cómo trabajé con Claude

Yo guié las decisiones de **arquitectura y negocio** (alcance del MVP, alma del producto, recortes, supuestos, score 2 factores, stack mínimo). Claude implementó **el código** (queries SQL, structs Go, componentes React, tests table-driven, Dockerfiles, nginx config). Ver `AI_LOG.md` para detalle de los prompts más relevantes.

---

## Iteración 2+ (si tuviera 2 semanas más)

1. Anticipación real (detección de cambio de comportamiento con histórico de pagos).
2. Envío de comunicaciones discriminadas desde la app, con plantillas por segmento.
3. Override del segmento por la analista con motivo obligatorio (lo que corté en decisión 4).
4. Vista CEO con KPIs ejecutivos (lo que corté en decisión 1).
5. Auth + multi-usuario con auditoría.
6. Importador del Sheet actual.
