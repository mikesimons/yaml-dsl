package withitems_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mikesimons/yaml-dsl/middleware"
	"github.com/mikesimons/yaml-dsl/middleware/withitems"
	"github.com/mikesimons/yaml-dsl/parser"
	"github.com/mikesimons/yaml-dsl/scripting/testparser"
	"github.com/mikesimons/yaml-dsl/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestHandler struct {
	types.ActionCommon `mapstructure:",squash"`
	Cb                 func(t *TestHandler)
}

func (t *TestHandler) Execute() (*types.ActionResult, error) {
	t.Cb(t)
	return &types.ActionResult{}, nil
}

func TestRun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WithItems middleware Suite")
}

func testFactory(handlerFn func(t *TestHandler)) (*withitems.Middleware, types.Action, map[string]interface{}, *parser.Dsl) {
	dsl := parser.New()
	dsl.ScriptParser = testparser.New()

	subject := &withitems.Middleware{Dsl: dsl}
	handler := &TestHandler{Cb: handlerFn}

	action := types.Action{
		Handler: func() types.Handler {
			return handler
		},
		Data: make(types.UnparsedAction),
	}

	vars := make(map[string]interface{})

	return subject, action, vars, dsl
}

var _ = Describe("WithItems middleware", func() {
	Context("with no item field present", func() {
		It("should execute action once", func() {
			executed := 0
			subject, action, vars, dsl := testFactory(func(t *TestHandler) {
				executed++
			})

			action.Data["name"] = "Test"

			subject.Execute(action, vars, middleware.Chain{DecodeFunc: dsl.Decode})

			Expect(executed).Should(Equal(1))
		})
	})

	Context("with key options unset", func() {
		It("should iterate a list specified in `with_items` field", func() {
			items := []interface{}{"item1", "item2"}

			subject, action, vars, dsl := testFactory(func(t *TestHandler) {
				expected := fmt.Sprintf("Item: %s", items[0])
				Expect(t.ActionCommon.Name).Should(Equal(expected))
				items = items[1:]
			})

			action.Data["name"] = "Item: {{item}}"
			action.Data["with_items"] = items

			subject.Execute(action, vars, middleware.Chain{DecodeFunc: dsl.Decode})
		})

		It("should iterate the result of an expression in `with_items` field", func() {
			items := []interface{}{"item1", "item2"}

			subject, action, vars, dsl := testFactory(func(t *TestHandler) {
				expected := fmt.Sprintf("Item: %s", items[0])
				Expect(t.ActionCommon.Name).Should(Equal(expected))
				items = items[1:]
			})

			action.Data["name"] = "Item: {{item}}"
			action.Data["with_items"] = "item1,item2"

			subject.Execute(action, vars, middleware.Chain{DecodeFunc: dsl.Decode})
		})

		It("should set the current item in `item` var", func() {
			subject, action, vars, dsl := testFactory(func(t *TestHandler) {
				Expect(t.ActionCommon.Name).Should(Equal("item1item1item1"))
			})

			action.Data["name"] = "{{item}}{{item}}{{item}}"
			action.Data["with_items"] = "item1"

			subject.Execute(action, vars, middleware.Chain{DecodeFunc: dsl.Decode})
		})
	})

	Context("with key options set", func() {
		It("should iterate a list specified in `my_items` field", func() {
			items := []interface{}{"item1", "item2"}

			subject, action, vars, dsl := testFactory(func(t *TestHandler) {
				expected := fmt.Sprintf("Item: %s", items[0])
				Expect(t.ActionCommon.Name).Should(Equal(expected))
				items = items[1:]
			})

			(*withitems.Middleware)(subject).Key = "my_items"
			(*withitems.Middleware)(subject).ItemKey = "my_item"
			action.Data["name"] = "Item: {{my_item}}"
			action.Data["my_items"] = items

			subject.Execute(action, vars, middleware.Chain{DecodeFunc: dsl.Decode})
		})

		It("should iterate the result of an expression in `my_items` field", func() {
			items := []interface{}{"item1", "item2"}

			subject, action, vars, dsl := testFactory(func(t *TestHandler) {
				expected := fmt.Sprintf("Item: %s", items[0])
				Expect(t.ActionCommon.Name).Should(Equal(expected))
				items = items[1:]
			})

			(*withitems.Middleware)(subject).Key = "my_items"
			(*withitems.Middleware)(subject).ItemKey = "my_item"
			action.Data["name"] = "Item: {{my_item}}"
			action.Data["my_items"] = "item1,item2"

			subject.Execute(action, vars, middleware.Chain{DecodeFunc: dsl.Decode})
		})

		It("should set the current item in `my_item` var", func() {
			subject, action, vars, dsl := testFactory(func(t *TestHandler) {
				Expect(t.ActionCommon.Name).Should(Equal("item1item1item1"))
			})

			(*withitems.Middleware)(subject).Key = "my_items"
			(*withitems.Middleware)(subject).ItemKey = "my_item"
			action.Data["name"] = "{{my_item}}{{my_item}}{{my_item}}"
			action.Data["my_items"] = "item1"

			subject.Execute(action, vars, middleware.Chain{DecodeFunc: dsl.Decode})
		})
	})

	Context("with bad data", func() {
		Context("(with_items expression)", func() {
			It("should return error", func() {
				subject, action, vars, dsl := testFactory(func(t *TestHandler) {})

				action.Data["name"] = "Test"
				action.Data["with_items"] = "BAD"
				dsl.ScriptParser.(*testparser.TestScriptParser).ParseExpressionError = errors.New("error")

				_, err := subject.Execute(action, vars, middleware.Chain{DecodeFunc: dsl.Decode})

				Expect(err.Error()).Should(Equal("error parsing expression: BAD: error"))
			})
		})

		Context("(expression embedded in item)", func() {
			It("should return error", func() {
				subject, action, vars, dsl := testFactory(func(t *TestHandler) {})

				action.Data["name"] = "Test"
				action.Data["with_items"] = []interface{}{"BAD"}
				dsl.ScriptParser.(*testparser.TestScriptParser).ParseEmbeddedError = errors.New("error")

				_, err := subject.Execute(action, vars, middleware.Chain{DecodeFunc: dsl.Decode})

				Expect(err.Error()).Should(Equal("error parsing expression: BAD: error"))
			})
		})
	})

	// TODO: Implement ActionResult handling
})
