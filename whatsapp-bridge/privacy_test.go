package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoMessageTextInLogs checks the bridge sources: no logging call receives the text of a message
func TestNoMessageTextInLogs(t *testing.T) {
	logFuncs := map[string]bool{"Print": true, "Printf": true, "Println": true, "Debugf": true, "Infof": true,
		"Warnf": true, "Errorf": true, "logf": true}
	textNames := map[string]bool{"content": true, "text": true, "Content": true, "Text": true}

	files, _ := filepath.Glob("*.go")
	checked := 0
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || !logFuncs[sel.Sel.Name] {
				return true
			}
			checked++
			for _, arg := range call.Args {
				ast.Inspect(arg, func(a ast.Node) bool {
					switch v := a.(type) {
					case *ast.Ident:
						if textNames[v.Name] {
							t.Errorf("%s: %s receives %q", fset.Position(call.Pos()), sel.Sel.Name, v.Name)
						}
					case *ast.CallExpr:
						// messageText(...), GetConversation(), GetText()...
						if s, ok := v.Fun.(*ast.SelectorExpr); ok && (s.Sel.Name == "GetConversation" || s.Sel.Name == "GetText") {
							t.Errorf("%s: %s receives a message text", fset.Position(call.Pos()), sel.Sel.Name)
						}
						if id, ok := v.Fun.(*ast.Ident); ok && id.Name == "messageText" {
							t.Errorf("%s: %s receives a message text", fset.Position(call.Pos()), sel.Sel.Name)
						}
					}
					return true
				})
			}
			return true
		})
	}
	if checked < 20 {
		t.Fatalf("only %d logging calls found: the check does not see the sources", checked)
	}
}

func TestChatAllowlistStrict(t *testing.T) {
	for _, raw := range []string{"", "  "} {
		t.Setenv("WHATSAPP_ALLOWED_CHATS", raw)
		if a, err := loadChatAllowlist(); a != nil || err != nil {
			t.Errorf("%q: an unset list allows every chat, got %v, %v", raw, a, err)
		}
	}
	for _, raw := range []string{" , ,", "15550000001, not a jid@@", "notajid@", "15550000001,@s.whatsapp.net"} {
		t.Setenv("WHATSAPP_ALLOWED_CHATS", raw)
		if a, err := loadChatAllowlist(); a != nil || err == nil {
			t.Errorf("%q: must be rejected, got %v, %v", raw, a, err)
		} else if strings.Contains(err.Error(), "1555") {
			t.Errorf("the error must not contain the entries: %v", err)
		}
	}
	t.Setenv("WHATSAPP_ALLOWED_CHATS", "+1 555 000 0001, 120363000000000001@g.us,")
	a, err := loadChatAllowlist()
	if err != nil || a == nil || len(a.jids) != 2 {
		t.Fatalf("valid list rejected: %v, %v", a, err)
	}
}
