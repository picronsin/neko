// Command protocolgen generates the Go and TypeScript surfaces shared by the
// WebSocket protocol. The JSON files under /protocol are the source of truth.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type eventSpec struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type errorSpec struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Retryable bool   `json:"retryable"`
}

type payloadSchema struct {
	Defs            map[string]schemaDefinition `json:"$defs"`
	AllOf           []payloadRule               `json:"allOf"`
	NoPayloadEvents []string                    `json:"x-no-payload-events"`
}

type payloadRule struct {
	If struct {
		Properties struct {
			Event struct {
				Const string   `json:"const"`
				Enum  []string `json:"enum"`
			} `json:"event"`
		} `json:"properties"`
	} `json:"if"`
}

type schemaDefinition struct {
	Type                 string                    `json:"type"`
	Required             []string                  `json:"required"`
	Properties           map[string]schemaProperty `json:"properties"`
	AdditionalProperties json.RawMessage           `json:"additionalProperties"`
}

type schemaProperty struct {
	Ref                  string          `json:"$ref"`
	Type                 json.RawMessage `json:"type"`
	Items                *schemaProperty `json:"items"`
	AdditionalProperties json.RawMessage `json:"additionalProperties"`
}

func main() {
	rootFlag := flag.String("root", "", "repository root; discovered from the current directory when omitted")
	flag.Parse()

	root, err := repositoryRoot(*rootFlag)
	if err != nil {
		fatal(err)
	}

	events := readJSON[[]eventSpec](filepath.Join(root, "protocol", "events.json"))
	errors := readJSON[[]errorSpec](filepath.Join(root, "protocol", "errors.json"))
	payloads := readJSON[payloadSchema](filepath.Join(root, "protocol", "payloads.schema.json"))
	validate(events, errors, payloads)

	goOutput, err := format.Source([]byte(generateGo(events, errors)))
	if err != nil {
		fatal(fmt.Errorf("format generated Go: %w", err))
	}
	write(filepath.Join(root, "server", "pkg", "protocol", "events_generated.go"), goOutput)

	goPayloadOutput, err := format.Source([]byte(generateGoPayloads(payloads)))
	if err != nil {
		fatal(fmt.Errorf("format generated Go payloads: %w", err))
	}
	write(filepath.Join(root, "server", "pkg", "protocol", "payloads_generated.go"), goPayloadOutput)

	tsOutput := generateTypeScript(events, errors)
	write(filepath.Join(root, "client", "src", "protocol", "events.generated.ts"), []byte(tsOutput))
	write(filepath.Join(root, "client", "src", "protocol", "payloads.generated.ts"), []byte(generatePayloadTypes(payloads)))
}

func repositoryRoot(explicit string) (string, error) {
	if explicit != "" {
		root, err := filepath.Abs(explicit)
		if err != nil {
			return "", err
		}
		return root, nil
	}

	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(directory, "protocol", "events.json")); err == nil {
			return directory, nil
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("unable to find repository root from %s", directory)
		}
		directory = parent
	}
}

func readJSON[T any](path string) T {
	data, err := os.ReadFile(path)
	if err != nil {
		fatal(err)
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		fatal(fmt.Errorf("parse %s: %w", path, err))
	}
	return result
}

func validate(events []eventSpec, errors []errorSpec, payloads payloadSchema) {
	seen := make(map[string]struct{}, len(events))
	for _, event := range events {
		if event.Name == "" || event.Value == "" {
			fatal(fmt.Errorf("event name and value are required"))
		}
		if _, ok := seen[event.Value]; ok {
			fatal(fmt.Errorf("duplicate event value %q", event.Value))
		}
		seen[event.Value] = struct{}{}
	}
	seen = make(map[string]struct{}, len(errors))
	for _, protocolError := range errors {
		if protocolError.Name == "" || protocolError.Value == "" {
			fatal(fmt.Errorf("error name and value are required"))
		}
		if _, ok := seen[protocolError.Value]; ok {
			fatal(fmt.Errorf("duplicate error value %q", protocolError.Value))
		}
		seen[protocolError.Value] = struct{}{}
	}

	knownEvents := make(map[string]struct{}, len(events))
	for _, event := range events {
		knownEvents[event.Value] = struct{}{}
	}
	coveredPayloads := make(map[string]struct{})
	for _, event := range payloads.NoPayloadEvents {
		if _, ok := knownEvents[event]; !ok {
			fatal(fmt.Errorf("payload schema marks unknown no-payload event %q", event))
		}
		if _, duplicate := coveredPayloads[event]; duplicate {
			fatal(fmt.Errorf("payload schema maps event %q more than once", event))
		}
		coveredPayloads[event] = struct{}{}
	}
	for _, rule := range payloads.AllOf {
		event := rule.If.Properties.Event
		values := append([]string{}, event.Enum...)
		if event.Const != "" {
			values = append(values, event.Const)
		}
		if len(values) == 0 {
			fatal(fmt.Errorf("payload schema contains a rule without an event"))
		}
		for _, value := range values {
			if _, ok := knownEvents[value]; !ok {
				fatal(fmt.Errorf("payload schema maps unknown event %q", value))
			}
			if _, duplicate := coveredPayloads[value]; duplicate {
				fatal(fmt.Errorf("payload schema maps event %q more than once", value))
			}
			coveredPayloads[value] = struct{}{}
		}
	}
	for event := range knownEvents {
		if _, ok := coveredPayloads[event]; !ok {
			fatal(fmt.Errorf("payload schema has no contract for event %q", event))
		}
	}
}

func generateGo(events []eventSpec, errors []errorSpec) string {
	var builder strings.Builder
	builder.WriteString("// Code generated by cmd/protocolgen; DO NOT EDIT.\n")
	builder.WriteString("package protocol\n\n")
	builder.WriteString("const (\n")
	for _, event := range events {
		fmt.Fprintf(&builder, "\t%s = %q\n", pascal(event.Name), event.Value)
	}
	builder.WriteString(")\n\n")
	builder.WriteString("var knownEvents = map[string]struct{}{\n")
	for _, event := range events {
		fmt.Fprintf(&builder, "\t%s: {},\n", pascal(event.Name))
	}
	builder.WriteString("}\n\nfunc IsKnownEvent(event string) bool {\n\t_, ok := knownEvents[event]\n\treturn ok\n}\n\n")
	builder.WriteString("type ErrorCode string\n\nconst (\n")
	for _, protocolError := range errors {
		fmt.Fprintf(&builder, "\t%s ErrorCode = %q\n", pascal(protocolError.Name), protocolError.Value)
	}
	builder.WriteString(")\n\nvar retryableErrors = map[ErrorCode]bool{\n")
	for _, protocolError := range errors {
		fmt.Fprintf(&builder, "\t%s: %t,\n", pascal(protocolError.Name), protocolError.Retryable)
	}
	builder.WriteString("}\n\nfunc IsRetryable(code ErrorCode) bool {\n\treturn retryableErrors[code]\n}\n")
	return strings.TrimRight(builder.String(), "\n") + "\n"
}

func generateGoPayloads(document payloadSchema) string {
	var builder strings.Builder
	builder.WriteString("// Code generated by cmd/protocolgen; DO NOT EDIT.\n")
	builder.WriteString("package protocol\n\n")
	builder.WriteString("// Protocol payload contracts generated from protocol/payloads.schema.json.\n")
	names := make([]string, 0, len(document.Defs))
	for name := range document.Defs {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		definition := document.Defs[name]
		if definition.Type != "object" {
			continue
		}
		goName := goPayloadTypeName(name)
		required := make(map[string]struct{}, len(definition.Required))
		for _, field := range definition.Required {
			required[field] = struct{}{}
		}
		fmt.Fprintf(&builder, "type %s struct {\n", goName)
		fields := make([]string, 0, len(definition.Properties))
		for field := range definition.Properties {
			fields = append(fields, field)
		}
		sort.Strings(fields)
		for _, field := range fields {
			property := definition.Properties[field]
			fieldType := goSchemaType(property)
			if name == "Error" && field == "code" {
				fieldType = "ErrorCode"
			}
			_, isRequired := required[field]
			if !isRequired && fieldType != "any" {
				fieldType = "*" + fieldType
			}
			tag := field
			if !isRequired {
				tag += ",omitempty"
			}
			fmt.Fprintf(&builder, "\t%s %s `json:%q`\n", goFieldName(field), fieldType, tag)
		}
		builder.WriteString("}\n\n")
	}
	return builder.String()
}

func goPayloadTypeName(name string) string {
	if name == "Error" {
		return "ProtocolErrorPayload"
	}
	return name + "Payload"
}

func goSchemaType(property schemaProperty) string {
	if property.Ref != "" {
		return goPayloadTypeName(strings.TrimPrefix(property.Ref, "#/$defs/"))
	}
	if len(property.Type) == 0 {
		return "any"
	}
	var typeName string
	if json.Unmarshal(property.Type, &typeName) == nil {
		return goPrimitiveSchemaType(typeName, property)
	}
	var typeNames []string
	if json.Unmarshal(property.Type, &typeNames) == nil {
		for _, name := range typeNames {
			if name == "null" {
				return "any"
			}
		}
		if len(typeNames) == 1 {
			return goPrimitiveSchemaType(typeNames[0], property)
		}
		return "any"
	}
	return "any"
}

func goPrimitiveSchemaType(typeName string, property schemaProperty) string {
	switch typeName {
	case "string":
		return "string"
	case "integer":
		return "int64"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "array":
		if property.Items == nil {
			return "[]any"
		}
		return "[]" + goSchemaType(*property.Items)
	case "object":
		if len(property.AdditionalProperties) == 0 {
			return "map[string]any"
		}
		var additional bool
		if json.Unmarshal(property.AdditionalProperties, &additional) == nil {
			return "map[string]any"
		}
		var additionalSchema schemaProperty
		if json.Unmarshal(property.AdditionalProperties, &additionalSchema) == nil {
			return "map[string]" + goSchemaType(additionalSchema)
		}
		return "map[string]any"
	default:
		return "any"
	}
}

func goFieldName(field string) string {
	parts := strings.FieldsFunc(field, func(r rune) bool { return r == '_' || r == '-' })
	var builder strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		builder.WriteString(strings.ToUpper(part[:1]))
		builder.WriteString(part[1:])
	}
	if builder.Len() == 0 {
		return "Value"
	}
	return builder.String()
}

func generateTypeScript(events []eventSpec, errors []errorSpec) string {
	var builder strings.Builder
	builder.WriteString("// Code generated by server/cmd/protocolgen; DO NOT EDIT.\n\n")
	builder.WriteString("export const PROTOCOL_EVENT = {\n")
	for _, event := range events {
		fmt.Fprintf(&builder, "  %s: %s,\n", event.Name, tsQuote(event.Value))
	}
	builder.WriteString("} as const\n\n")
	builder.WriteString("export type ProtocolEvent = (typeof PROTOCOL_EVENT)[keyof typeof PROTOCOL_EVENT]\n\n")
	builder.WriteString("export const PROTOCOL_ERROR = {\n")
	for _, protocolError := range errors {
		fmt.Fprintf(&builder, "  %s: %s,\n", protocolError.Name, tsQuote(protocolError.Value))
	}
	builder.WriteString("} as const\n\n")
	builder.WriteString("export type ProtocolErrorCode = (typeof PROTOCOL_ERROR)[keyof typeof PROTOCOL_ERROR]\n\n")
	builder.WriteString("export const RETRYABLE_PROTOCOL_ERRORS: ReadonlySet<ProtocolErrorCode> = new Set([\n")
	for _, protocolError := range errors {
		if protocolError.Retryable {
			fmt.Fprintf(&builder, "  PROTOCOL_ERROR.%s,\n", protocolError.Name)
		}
	}
	builder.WriteString("])\n")
	return builder.String()
}

func generatePayloadTypes(document payloadSchema) string {
	var builder strings.Builder
	builder.WriteString("// Code generated by server/cmd/protocolgen; DO NOT EDIT.\n\n")
	builder.WriteString("import type { ProtocolErrorCode } from './events.generated'\n\n")
	names := make([]string, 0, len(document.Defs))
	for name := range document.Defs {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		definition := document.Defs[name]
		if definition.Type != "object" {
			continue
		}
		required := make(map[string]struct{}, len(definition.Required))
		for _, field := range definition.Required {
			required[field] = struct{}{}
		}
		fmt.Fprintf(&builder, "export interface %sPayload {\n", name)
		fields := make([]string, 0, len(definition.Properties))
		for field := range definition.Properties {
			fields = append(fields, field)
		}
		sort.Strings(fields)
		for _, field := range fields {
			property := definition.Properties[field]
			_, isRequired := required[field]
			optional := "?"
			if isRequired {
				optional = ""
			}
			fieldType := schemaType(property)
			if name == "Error" && field == "code" {
				fieldType = "ProtocolErrorCode"
			}
			fmt.Fprintf(&builder, "  %s%s: %s\n", field, optional, fieldType)
		}
		builder.WriteString("}\n\n")
	}
	return strings.TrimRight(builder.String(), "\n") + "\n"
}

func schemaType(property schemaProperty) string {
	if property.Ref != "" {
		name := strings.TrimPrefix(property.Ref, "#/$defs/")
		return name + "Payload"
	}
	if len(property.Type) == 0 {
		return "unknown"
	}
	var typeName string
	if json.Unmarshal(property.Type, &typeName) == nil {
		return primitiveSchemaType(typeName, property)
	}
	var typeNames []string
	if json.Unmarshal(property.Type, &typeNames) == nil {
		parts := make([]string, 0, len(typeNames))
		for _, name := range typeNames {
			parts = append(parts, primitiveSchemaType(name, property))
		}
		return strings.Join(parts, " | ")
	}
	return "unknown"
}

func primitiveSchemaType(typeName string, property schemaProperty) string {
	switch typeName {
	case "string":
		return "string"
	case "integer", "number":
		return "number"
	case "boolean":
		return "boolean"
	case "null":
		return "null"
	case "array":
		if property.Items == nil {
			return "unknown[]"
		}
		return "Array<" + schemaType(*property.Items) + ">"
	case "object":
		if len(property.AdditionalProperties) == 0 {
			return "Record<string, unknown>"
		}
		var additional bool
		if json.Unmarshal(property.AdditionalProperties, &additional) == nil {
			if additional {
				return "Record<string, unknown>"
			}
			return "Record<string, never>"
		}
		var additionalSchema schemaProperty
		if json.Unmarshal(property.AdditionalProperties, &additionalSchema) == nil {
			return "Record<string, " + schemaType(additionalSchema) + ">"
		}
		return "Record<string, unknown>"
	default:
		return "unknown"
	}
}

func tsQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "\\'") + "'"
}

func pascal(name string) string {
	parts := strings.Split(strings.ToLower(name), "_")
	for index, part := range parts {
		if part == "" {
			continue
		}
		parts[index] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

func write(path string, data []byte) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
