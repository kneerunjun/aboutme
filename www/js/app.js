(function() {
    angular.module("aboutme", ["ngRoute"]).config(function($routeProvider, $interpolateProvider, $provide, $locationProvider) {
        // with GO Lang frameworks this can help to have angular a distinct space 
        $interpolateProvider.startSymbol("{[")
        $interpolateProvider.endSymbol("]}")
        $locationProvider.html5Mode({
            enabled: true,
            requireBase: true
        });
        // For now we are just blocking the routes - but we would certainly want them re-enabled when we have multiple pages
        $routeProvider
            .when("/", {
                templateUrl: "/views/splash.html"
            })
        //     .when("/products", {
        //         templateUrl: "/views/products-list.html"
        //     })
        //     .when("/blogs", {
        //         templateUrl: "/views/blogs-list.html"
        //     })
        //     .when("/blogs/:id", {
        //         templateUrl: "/views/blogs-read.html"
        //     })
        //     .when("/about", {
        //         templateUrl: "/views/about.html"
        //     })
        //     .when("/products/:id", {
        //         templateUrl: "/views/product-detail.html"
        //     })
        //     .when("/testpay", {
        //         templateUrl: "/views/test-pay.html"
        //     })
        $provide.provider("emailPattern", function() {
            this.$get = function() {
                // [\w] is the same as [A-Za-z0-9_-]
                // 3 groups , id, provider , domain also a '.' in between separated by @
                // we are enforcing a valid email id 
                // email id can have .,_,- in it and nothing more 
                return /^[\w-._]+@[\w]+\.[a-z]+$/
            }
        })
        $provide.provider("passwdPattern", function() {
            this.$get = function() {
                // here for the password the special characters that are not allowed are being singled out and denied.
                // apart form this all the characters will be allowed
                // password also has a restriction on the number of characters in there
                return /^[\w-!@#%&?_]{8,16}$/
            }
        })
    })
})()