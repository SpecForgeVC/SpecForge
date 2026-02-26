# SpecForge Iterative Project Cataloguing Protocol (IPCP) - Advanced Guide

This document defines the deterministic protocol for performing a full project import into SpecForge. Follow these instructions strictly to achieve 100% completeness score.

---

## 🚀 Execution Strategy: The "Deep Scan" Protocol

### Phase 1: Contextual Research
- **Scan Root**: Identify tech stack, entry points, and project structure.
- **Dependency Analysis**: Identify core internal modules and external integrations.
- **API Discovery**: Search for controllers, route definitions, and documentation (OpenAPI, etc.).

### Phase 2: Local Scaffolding (Persistence)
- Create a `.specforge/` directory.
- Create local `.json` files for each category. **You must keep these local files updated.** If the server requests a rescan, update the local file before resubmitting.

### Phase 3: Incremental Submission
- Do not attempt to submit the entire project in one giant tool call if the project is large.
- Submit categories as they are ready. The server will merge snapshots keyed by `project_id`.

### Phase 4: Self-Assessment & Refinement
- Read the `completeness_score` in the tool response.
- Review `missing_categories` and `self_assessment_prompt`.
- Perform targeted searches in the codebase for the specifically missing items.

---

## 📋 Comprehensive Category Requirements

### 1. `project_overview` (One Object)
**Purpose**: High-level identity.
- `name`: Target project name.
- `description`: Detailed technical summary.
- `domain`: Industry or software domain (e.g., "Developer Tools", "E-commerce").
- `primary_language`: Main programming language.
- `architecture_pattern`: E.g., "Microservices", "Monolith", "Layered".

### 2. `tech_stack` (One Object)
**Purpose**: Technology fingerprint.
- `languages`: Array of all languages used.
- `frameworks`: Web/ORM/Testing frameworks.
- `databases`: SQL/NoSQL systems.
- `infrastructure`: Cloud providers, Docker, K8s.
- `build_tools`: npm, webpack, go build, etc.

### 3. `modules` (Array of Objects)
**Purpose**: Logical code boundaries.
- `name`: Unique identifier for the module.
- `description`: What it does.
- `responsibilities`: List of key duties.
- `risk_level`: `LOW` | `MEDIUM` | `HIGH`.
- `change_sensitivity`: `LOW` | `MEDIUM` | `HIGH`.
- `dependencies`: Names of other modules this one relies on.

### 4. `apis` (Array of Objects)
**Purpose**: External and internal communication surfaces.
- `endpoint`: URL path or identifier.
- `method`: `GET` | `POST` | `PUT` | `PATCH` | `DELETE`.
- `auth_type`: `Bearer`, `Session`, `None`.
- `request_schema` / `response_schema`: (Optional) JSON schema objects.

### 5. `data_models` (Array of Objects)
**Purpose**: Core entities and relationships.
- `name`: Entity name (e.g., "User").
- `relationships`: Array containing `{ target: "EntityName", type: "1:N" }`.

### 6. `contracts` (Array of Objects)
**Purpose**: Formal interfaces and stability.
- `name`: Interface/Function name.
- `contract_type`: `REST` | `GRAPHQL` | `EVENT` | `INTERNAL_FUNCTION`.
- `source_module`: The module defining this contract.
- `stability_score`: 0.0 to 1.0.

### 7. `risks` (Array of Objects)
**Purpose**: Architectural and security concerns.
- `area`: Component at risk.
- `severity`: `LOW` | `MEDIUM` | `HIGH` | `CRITICAL`.
- `description`: Why it's a risk.

### 8. `change_sensitivity` (Array of Objects)
**Purpose**: Mapping of where changes cause the most impact.
- `module`: Target module.
- `sensitivity`: `LOW` | `MEDIUM` | `HIGH`.
- `reason`: Explanation.

---

## 📐 Data Structure Examples

### Example: `modules` + `contracts` Alignment
Ensure that `source_module` in a contract refers to a name defined in the `modules` category.

```json
// .specforge/modules.json
[
  { "name": "AuthService", "description": "Handles JWT...", "risk_level": "HIGH" }
]

// .specforge/contracts.json
[
  { "name": "ValidateToken", "contract_type": "INTERNAL_FUNCTION", "source_module": "AuthService" }
]
```

---

## 🛠 Tool Usage Checklist

### `mcp_specforge_init_project_import`
Call this first to get your `SessionID`.
- Provide accurate `language_stack` and `repository_type`.

### `mcp_specforge_submit_project_snapshot`
Call this repeatedly to build the profile.
- `project_id`: Must match the initial call.
- `snapshot_version`: Use a timestamp or increment (e.g., `v1.0.1`).
- `snapshot_payload`: Construct this by reading your `.specforge/*.json` files.
- `final_submission`: Set to `true` when the `completeness_score` is high to indicate a major milestone.

### `mcp_specforge_finalize_project_import`
Call this tool ONLY when the project is 100% catalogued and the server has returned a score >= 95.
- `project_id`: The identifier for the project.
- **Effect**: This officially closes the import session and transitions the user's UI to the project dashboard.

---

## ⚠️ Finality Criteria
The import is complete only when:
1. All 8 core categories are populated.
2. The `completeness_score` is >= 95.
3. You have called `mcp_specforge_finalize_project_import`.
4. The tool response returns `status: finalized`.
