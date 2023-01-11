import "./index.css"

import { ChakraProvider, Heading } from "@chakra-ui/react"
import { Login, Private, Public, SignUp } from "./routes/"
import { RouterProvider, createHashRouter } from "react-router-dom"

import App from "./App"
import { Auth0Provider } from "@auth0/auth0-react"
import React from "react"
import { createRoot } from "react-dom/client"
import { getConfig } from "./settings/config"
import history from "./utils/history"

const container = document.getElementById("root")
const root = createRoot(container!) // createRoot(container!) if you use TypeScript

const router = createHashRouter([
  {
    path: "/",
    element: <App />,
    errorElement: <Heading>Page not Found</Heading>,
    children: [
      {
        path: "/",
        element: <Heading>Landing Page</Heading>,
      },
      {
        path: "/private/data",
        element: <Private />,
      },
      {
        path: "/public/data",
        element: <Public />,
      },
      {
        path: "/login",
        element: <Login />,
      },
      {
        path: "/register",
        element: <SignUp />,
      },
    ],
  },
])

const onRedirectCallback = (appState: any) => {
  history.push(
    appState && appState.returnTo ? appState.returnTo : window.location.pathname
  )
}

const config = getConfig()

const providerConfig = {
  domain: config.domain,
  clientId: config.clientId,
  ...(config.audience ? { audience: config.audience } : null),
  redirectUri: window.location.origin,
  onRedirectCallback,
}

root.render(
  <React.StrictMode>
    <Auth0Provider {...providerConfig}>
      <ChakraProvider>
        <RouterProvider router={router} />
      </ChakraProvider>
    </Auth0Provider>
    ,
  </React.StrictMode>
)
