package lib

import (
	"net/http"
)

type ReqHookFunc = func(req *http.Request, r Redirect);
type ResHookFunc = func(req *http.Response, r Redirect);

// Interface defining all functions a plugin should have in this app
// Plugins may implement a simple `return` to do nothing
// Plugins have to expose a `Plugin CustomsPlugin` symbol to expose the functions
type CustomsPlugin interface {
    // Add more configuration flags. Use the flag package default
    AddConfigFlag();

    // Called before the request is made. Where request can be altered
    PreRequestHook(req *http.Request, r Redirect);

    // Called after the request is made. Where response can be altered.
    PostRequestHook(res *http.Response, r Redirect);
}
