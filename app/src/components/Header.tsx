import { Avatar, Box, Button, Flex, HStack, Text } from "@chakra-ui/react"

import { NavLink } from "react-router-dom"
import { useAuth0 } from "@auth0/auth0-react"
import { useAuthContext } from "../App"

// import { useAuthContext } from "../App"

export function Header() {
  const { user, isAuthenticated, loginWithRedirect, logout, getIdTokenClaims } =
    useAuth0()
  const { user: currentUser, signout } = useAuthContext()
  console.log("currentUser: ", currentUser)
  // const { user, signout } = useAuthContext()

  const logoutWithRedirect = () =>
    logout({
      returnTo: window.location.origin,
    })
  console.log("user: ", user)

  return (
    <Flex justify="space-between" align="center" bgColor="gray.200" p="4">
      <Text fontSize="20"></Text>
      <HStack as="nav" spacing="4" p="4" bgColor="gray.200" borderRadius="8">
        <NavLink to="/">
          <Button colorScheme="blue" variant="outline">
            Home
          </Button>
        </NavLink>
        <NavLink to="/private/data">
          <Button colorScheme="blue" variant="outline">
            Private Data
          </Button>
        </NavLink>
        <NavLink to="/public/data">
          <Button colorScheme="blue" variant="outline">
            Public Data
          </Button>
        </NavLink>
        {isAuthenticated ? (
          <>
            <Button
              colorScheme="blue"
              variant="ghost"
              onClick={logoutWithRedirect}
            >
              Log out
            </Button>
          </>
        ) : (
          <Button colorScheme="blue" onClick={loginWithRedirect}>
            Login with Auth0
          </Button>
        )}
        {currentUser?.isAuthenticated ? (
          <Button
            colorScheme="blue"
            variant="ghost"
            onClick={() => signout?.()}
          >
            Log out
          </Button>
        ) : (
          <NavLink to="/login">
            <Button colorScheme="pink">Login Form</Button>
          </NavLink>
        )}
      </HStack>
      <Box className="user-info">
        {isAuthenticated && (
          <Avatar name={user?.nickname} src={user?.picture} />
        )}
        {/* <h6 className="d-inline-block">{user?.name}</h6> */}
      </Box>
    </Flex>
  )
}
