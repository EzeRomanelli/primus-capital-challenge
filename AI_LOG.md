# AI_LOG — Uso de Claude Code

> **El principio operativo:** yo guié las decisiones de **arquitectura y negocio**, Claude implementó **el código**.
>
> Cuando Claude propuso algo que contradecía el alcance, lo corregí explícitamente. Cuando aportó valor, lo aproveché.

---

## Lo que yo decidí (sin delegar)

- **Alcance del MVP en 3 días.** Cortar anticipación, mantener foco + gestión mínima. Recortar a 1 flujo principal sólido.
- **Score con 2 factores, no 4.** Claude propuso 4 (urgencia + impacto + descuido + riesgo); yo lo bajé porque la complejidad inicial sin filtrar contra el alcance es el patrón típico de la IA.
- **Stack: Go + Chi + pgx + SQL plano.** Claude propuso `sqlc` (codegen). Yo cuestioné si necesitaba un ORM siquiera; 4 tablas no lo justifican.
- **Premisa por escrito al Day 0: "simplicidad para apagar incendios".** Esa frase fue mi filtro defensivo en cada decisión posterior.
- **No contactar a `jrain@primuscapital.cl`.** Decisión consciente, documentada en `DECISIONS.md` #6.
- **Cortar override del segmento.** El doc pide UN flujo principal end-to-end; dos a medias es peor.
- **Cortar la vista CEO.** El doc dice literal "2 flujos sólidos > 3 a medias".

## Lo que Claude implementó (con guía mía)

- Queries SQL (joins, agregados de facturas pendientes, índices con WHERE).
- Structs Go con json tags + scoring puro (función testeable con tabla de casos).
- Suggester rule-based + tests table-driven.
- Recalc fire-and-forget con context propio (no el del request).
- Endpoints HTTP con validación de inputs y errores `{error, code}` uniformes.
- Componentes React (Dashboard, ClienteDetail, GestionDialog), hooks de TanStack Query con `invalidateQueries`.
- Dockerfiles multi-stage + nginx.conf con SPA fallback + docker-compose con healthchecks encadenados.
- Diagrama ER en graphviz (.dot) y DBML.

## Dónde corregí a Claude

1. **Plan de 4426 líneas → 250.** Cuando le pedí el plan inicial sobre-detalló siguiendo su skill de planning sin filtrar contra el alcance. Hice cumplir el principio de simplicidad, lo redujo a 250 líneas.
2. **Score 4 factores → 2.** Filtré contra el alcance.
3. **sqlc → pgx + SQL plano.** Filtré contra el alcance.
4. **Tolerancia Corporativo 75 → 30.** Claude auto-detectó la inconsistencia entre `tolerancia_dias` (operacional, por segmento) y `payment_terms_dias` (contractual, por cliente) en una auto-revisión. Yo acepté la corrección tras entender la diferencia.
5. **Sobre-construcción día 3.** Después de un primer round terminé con 8 decisiones grandes + 8 chicas + tooltip + filtros + 5 personajes hardcoded + mockup Pencil + PITCH.md. Auditando contra el doc, **corté todo lo que estaba sobre el mínimo** y rehice el repo desde cero con `git init` para que el historial cuente la versión defendible, no la sobre-construida.

## Patrón que noto en mi uso de IA

- **Claude default: más detalle, más factores, más capas.** Útil para boilerplate, peligroso para alcance.
- **Mi rol: filtrar contra el alcance.** Cada vez que Claude propuso algo, le aplicaba la pregunta "¿esto resuelve el dolor principal o es lindo de tener?".
- **Por escrito ayuda más que en mi cabeza.** La premisa de simplicidad la escribí literal en `DECISIONS.md` y la cité varias veces ante tentaciones de scope creep.
