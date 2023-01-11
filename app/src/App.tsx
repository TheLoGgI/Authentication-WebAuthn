import "./App.css"

import { Header, Layout } from "./components"
import { createContext, useContext, useEffect, useState } from "react"

import { Outlet } from "react-router-dom"

const AppContext = createContext({
  feedbackBehavior: "bad",
  isDevelopment: true,
})

type User = {
  username: string
  uid: string
  firstName: string
  lastName: string
  email: string
  isAuthenticated: boolean
}

type AuthContextType = {
  user: User | null
  signIn: (() => void) | null
  signout: (() => void) | null
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  signIn: () => undefined,
  signout: () => undefined,
})

function useAuth(): AuthContextType {
  const [user, setUser] = useState<User | null>(null)

  const signIn = () => {
    const { exp, ...user } = JSON.parse(localStorage.getItem("user") ?? "{}")
    console.log("user: ", user)

    if (user) {
      setUser({
        ...user,
        isAuthenticated: new Date() < new Date(exp),
      })
    }
  }

  const signout = () => {
    setUser(null)
    localStorage.removeItem("user")
    sessionStorage.removeItem("user")
  }

  return {
    user,
    signIn,
    signout,
  }
}

export const useAuthContext = (): AuthContextType => {
  return useContext<AuthContextType>(AuthContext)
}

function App() {
  const auth = useAuth()
  return (
    <AuthContext.Provider value={auth}>
      <Header /* context={contextState} */ />
      <Layout>
        <Outlet context={AppContext} />
      </Layout>
    </AuthContext.Provider>
  )
}

export default App
