package scripting

import (
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"

	"postie/pkg/client"
	"postie/pkg/httprequest"
)

// Engine executes JavaScript response handler scripts
type Engine struct {
	vm      *goja.Runtime
	context *ScriptContext
	results *ScriptExecutionResult
}

// NewEngine creates a new JavaScript execution engine
func NewEngine(context *ScriptContext) *Engine {
	engine := &Engine{
		vm:      goja.New(),
		context: context,
		results: &ScriptExecutionResult{
			Tests:      make([]*TestResult, 0),
			Assertions: make([]*AssertionError, 0),
			Logs:       make([]string, 0),
			Globals:    make(map[string]interface{}),
		},
	}

	engine.setupClientAPI()
	engine.setupResponseObject()
	engine.setupRequestObject()
	engine.setupEnvironmentVariables()

	return engine
}

// Execute runs the JavaScript script and returns the results
func (e *Engine) Execute(script string) *ScriptExecutionResult {
	defer func() {
		if r := recover(); r != nil {
			e.results.Error = fmt.Errorf("script panic: %v", r)
		}
	}()

	_, err := e.vm.RunString(script)
	if err != nil {
		e.results.Error = fmt.Errorf("script execution error: %w", err)
	}

	// Copy globals back to context
	if e.context.Globals != nil {
		e.results.Globals = e.context.Globals.GetAll()
	}

	return e.results
}

// setupClientAPI sets up the client object with test, assert, log, and global methods
func (e *Engine) setupClientAPI() {
	client := e.vm.NewObject()

	// client.test(name, function)
	client.Set("test", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			e.results.Error = fmt.Errorf("client.test() requires 2 arguments: name and function")
			return goja.Undefined()
		}

		name := call.Argument(0).String()
		testFunc, ok := goja.AssertFunction(call.Argument(1))
		if !ok {
			e.results.Error = fmt.Errorf("client.test() second argument must be a function")
			return goja.Undefined()
		}

		result := &TestResult{
			Name:   name,
			Passed: true,
		}

		// Execute the test function
		_, err := testFunc(goja.Undefined())
		if err != nil {
			result.Passed = false
			result.Error = err.Error()
		}

		e.results.Tests = append(e.results.Tests, result)
		return goja.Undefined()
	})

	// client.assert(condition, message)
	client.Set("assert", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(e.vm.NewGoError(fmt.Errorf("client.assert() requires at least 1 argument")))
		}

		condition := call.Argument(0).ToBoolean()
		message := "Assertion failed"
		if len(call.Arguments) >= 2 {
			message = call.Argument(1).String()
		}

		if !condition {
			assertErr := &AssertionError{
				Message: message,
			}
			e.results.Assertions = append(e.results.Assertions, assertErr)
			panic(e.vm.NewGoError(assertErr))
		}

		return goja.Undefined()
	})

	// client.log(...messages)
	client.Set("log", func(call goja.FunctionCall) goja.Value {
		messages := make([]string, len(call.Arguments))
		for i, arg := range call.Arguments {
			messages[i] = arg.String()
		}

		logMessage := fmt.Sprint(messages)
		e.results.Logs = append(e.results.Logs, logMessage)
		return goja.Undefined()
	})

	// client.global object for global variables
	global := e.vm.NewObject()

	// client.global.set(name, value)
	global.Set("set", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			e.results.Error = fmt.Errorf("client.global.set() requires 2 arguments: name and value")
			return goja.Undefined()
		}

		name := call.Argument(0).String()
		value := call.Argument(1).Export()

		if e.context.Globals != nil {
			e.context.Globals.Set(name, value)
		}

		return goja.Undefined()
	})

	// client.global.get(name)
	global.Set("get", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}

		name := call.Argument(0).String()

		if e.context.Globals != nil {
			if value, exists := e.context.Globals.Get(name); exists {
				return e.vm.ToValue(value)
			}
		}

		return goja.Undefined()
	})

	// client.global.clear(name)
	global.Set("clear", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}

		name := call.Argument(0).String()

		if e.context.Globals != nil {
			e.context.Globals.Clear(name)
		}

		return goja.Undefined()
	})

	// client.global.isEmpty()
	global.Set("isEmpty", func(call goja.FunctionCall) goja.Value {
		if e.context.Globals == nil {
			return e.vm.ToValue(true)
		}

		globals := e.context.Globals.GetAll()
		return e.vm.ToValue(len(globals) == 0)
	})

	client.Set("global", global)

	e.vm.Set("client", client)
}

// setupResponseObject sets up the response object in the script context
func (e *Engine) setupResponseObject() {
	if e.context.Response == nil {
		return
	}

	response := e.vm.NewObject()

	// response.status
	response.Set("status", e.context.Response.StatusCode)

	// response.statusText
	response.Set("statusText", e.context.Response.Status)

	// response.headers
	headers := make(map[string]string)
	for key, values := range e.context.Response.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	response.Set("headers", headers)

	// response.body
	body, err := e.context.Response.Text()
	if err == nil {
		// Try to parse as JSON
		var jsonBody interface{}
		if err := json.Unmarshal([]byte(body), &jsonBody); err == nil {
			response.Set("body", jsonBody)
		} else {
			response.Set("body", body)
		}
	}

	// response.contentType
	response.Set("contentType", e.context.Response.ContentType())

	e.vm.Set("response", response)
}

// setupRequestObject sets up the request object in the script context
func (e *Engine) setupRequestObject() {
	if e.context.Request == nil {
		return
	}

	request := e.vm.NewObject()

	// request.method
	request.Set("method", e.context.Request.Method)

	// request.url
	if e.context.Request.URL != nil {
		request.Set("url", e.context.Request.URL.Raw)
	}

	// request.headers
	headers := make(map[string]string)
	for _, header := range e.context.Request.Headers {
		headers[header.Name] = header.Value
	}
	request.Set("headers", headers)

	e.vm.Set("request", request)
}

// setupEnvironmentVariables sets up environment variables in the script context
func (e *Engine) setupEnvironmentVariables() {
	if e.context.Env == nil {
		return
	}

	e.vm.Set("env", e.context.Env)
}

// ExecuteResponseHandler executes a response handler script
func ExecuteResponseHandler(handler *httprequest.ResponseHandler, response *client.Response, request *httprequest.Request, env map[string]interface{}, globals *GlobalStore) *ScriptExecutionResult {
	if handler == nil {
		return &ScriptExecutionResult{
			Tests:      make([]*TestResult, 0),
			Assertions: make([]*AssertionError, 0),
			Logs:       make([]string, 0),
			Globals:    make(map[string]interface{}),
		}
	}

	context := &ScriptContext{
		Request:  request,
		Response: response,
		Env:      env,
		Globals:  globals,
	}

	engine := NewEngine(context)

	// Execute inline script or load from file
	script := handler.Script
	if handler.Type == httprequest.HandlerTypeFile {
		// TODO: Load script from file
		return &ScriptExecutionResult{
			Error: fmt.Errorf("file-based response handlers not yet implemented"),
		}
	}

	return engine.Execute(script)
}
