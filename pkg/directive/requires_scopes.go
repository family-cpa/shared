package directive

import (
	"entgo.io/contrib/entgql"
	"github.com/vektah/gqlparser/v2/ast"
)

func RequiresScopes(scopes ...string) entgql.Directive {
	groups := make(ast.ChildValueList, 0, len(scopes))
	for _, scope := range scopes {
		groups = append(groups, &ast.ChildValue{
			Value: &ast.Value{
				Kind: ast.ListValue,
				Children: ast.ChildValueList{
					{Value: &ast.Value{Kind: ast.StringValue, Raw: scope}},
				},
			},
		})
	}
	return entgql.NewDirective("requiresScopes", &ast.Argument{
		Name:  "scopes",
		Value: &ast.Value{Kind: ast.ListValue, Children: groups},
	})
}

func Schema(s *ast.Schema) {
	if s == nil {
		return
	}
	for name, def := range s.Types {
		if def == nil || name == "Query" || name == "Mutation" {
			continue
		}
		for _, field := range def.Fields {
			if !has(field, "requiresScopes") || field.Type == nil {
				continue
			}
			field.Type.NonNull = false
		}
	}
}

func has(field *ast.FieldDefinition, name string) bool {
	if field == nil {
		return false
	}
	for _, d := range field.Directives {
		if d != nil && d.Name == name {
			return true
		}
	}
	return false
}
