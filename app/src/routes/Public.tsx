import {
  Avatar,
  Box,
  Divider,
  Flex,
  HStack,
  Heading,
  ListItem,
  Spinner,
  Text,
  UnorderedList,
} from "@chakra-ui/react"
import { useAuth0, withAuthenticationRequired } from "@auth0/auth0-react"

import { useAuthContext } from "../App"
import { useEffect } from "react"

// import { useAuthContext } from "../App"

export const Public = () => {
  return (
    <Box>
      <Heading>Public data</Heading>
      <Text>Always available</Text>
    </Box>
  )
}

export const Private = () => {
  // const { user } = useAuthContext()
  const { user, isAuthenticated } = useAuth0()
  const { user: currentUser, signIn } = useAuthContext()

  useEffect(() => {
    signIn?.()
  }, [])

  if (currentUser?.isAuthenticated)
    return (
      <Box>
        <Heading>Private data</Heading>

        <Box bg="purple.100" borderRadius="10" p="8" mt="4">
          <HStack spacing="4">
            {
              <Avatar
                name="random person"
                src="https://thispersondoesnotexist.com/image"
                size="xl"
              />
            }
            <Text fontWeight="bold" textTransform="capitalize">
              {currentUser.username}
            </Text>
          </HStack>
          <Divider my="4" />
          <UnorderedList>
            <ListItem>Uid: {currentUser?.uid}</ListItem>
            <ListItem>Email: {currentUser?.email}</ListItem>
            <ListItem>
              isAuthenticated: {String(currentUser?.isAuthenticated)}
            </ListItem>
          </UnorderedList>
        </Box>
      </Box>
    )

  if (isAuthenticated)
    return (
      <Box>
        <Heading>Private data</Heading>

        <Box bg="purple.100" borderRadius="10" p="8" mt="4">
          <HStack spacing="4">
            {<Avatar name={user?.nickname} src={user?.picture} />}
            <Text>
              {user?.name}, {user?.nickname}
            </Text>
          </HStack>
          <Divider my="8" />
          <UnorderedList>
            <ListItem>Email: {user?.email} </ListItem>
            <ListItem>Verified: {String(user?.email_verified)} </ListItem>
            {user?.updated_at && (
              <ListItem>
                <>
                  Updated At: {new Date(user?.updated_at).toLocaleDateString()}
                </>
              </ListItem>
            )}
          </UnorderedList>
        </Box>
      </Box>
    )

  return (
    <Box color="red.400">
      <Heading>Private data</Heading>
      <Text>Is not available to this user</Text>
    </Box>
  )
}

// export default withAuthenticationRequired(Private, {
//   onRedirecting: () => <Spinner />,
// })
