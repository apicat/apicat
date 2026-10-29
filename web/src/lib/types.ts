export interface Schema {
  Name: string
  Type: string
  Format: string
  Description: string
  Required: boolean
  Nullable: boolean
  Deprecated: boolean
  Enum: string[] | null
  Default: string
  Example: string
  JSONExample: string
  JSONDefault: string
  MinLength: number | null
  MaxLength: number | null
  Minimum: number | null
  Maximum: number | null
  MinItems: number | null
  MaxItems: number | null
  Pattern: string
  Properties: Schema[] | null
  Items: Schema | null
  AllOf: Schema[] | null
  OneOf: Schema[] | null
  AnyOf: Schema[] | null
  RefName: string
  IsRef: boolean
  Depth: number
  Recursive: boolean
  Truncated: boolean
}
export interface Parameter {
  Name: string
  In: string
  Description: string
  Required: boolean
  Deprecated: boolean
  Schema: Schema | null
  Example: string
  Style: string
  Explode: boolean | null
}
export interface Media {
  MediaType: string
  Schema: Schema | null
  Example: string
}
export interface RequestBody {
  Description: string
  Required: boolean
  Content: Media[] | null
}
export interface APIResponse {
  StatusCode: string
  Description: string
  Content: Media[] | null
}
export interface Server {
  URL: string
  Description: string
}
export interface Endpoint {
  Method: string
  Path: string
  Summary: string
  Description: string
  OperationID: string
  Deprecated: boolean
  Parameters: Parameter[] | null
  RequestBody: RequestBody | null
  Responses: APIResponse[] | null
  Anchor: string
  LintHints: string[] | null
  Servers: Server[] | null
}
export interface Component {
  Name: string
  Description: string
  Schema: Schema | null
  Parameters: Parameter[] | null
  RequestBody: RequestBody | null
  Response: APIResponse | null
  Definition: string
}
export interface APIDoc {
  Title: string
  Description: string
  Version: string
  Servers: Server[] | null
  Tags: { Name: string; Description: string }[] | null
  TagGroups: { Tag: { Name: string; Description: string }; Endpoints: Endpoint[] }[] | null
  Security: { Name: string; Type: string; Scheme: string }[] | null
  ComponentGroups: { Kind: string; Title: string; Definitions: Component[] }[] | null
  LintTotal: number
  LintCounts: { Rule: string; Count: number; Message: string }[] | null
}
export interface Pair {
  key: string
  value: string
  enabled: boolean
}
export interface Variable extends Pair {
  secret: boolean
}
export interface Environment {
  id: string
  name: string
  baseUrl: string
  variables: Variable[]
}
export interface Workspace {
  revision: number
  activeEnvironment: string
  variables: Variable[]
  environments: Environment[]
}
export interface Session {
  token: string
  specDir: string
  workspacePath: string
}
export interface RequestDraft {
  queryResolved?: boolean
  method: string
  url: string
  environmentId: string
  headers: Pair[]
  query: Pair[]
  path: Pair[]
  body: string
  timeoutMs: number
}
export interface RunResponse {
  status: number
  statusText: string
  headers: Record<string, string[]>
  body: string
  encoding: string
  durationMs: number
  size: number
  truncated: boolean
  url: string
}
