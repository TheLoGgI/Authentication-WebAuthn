import {
  Box,
  Button,
  Input as ChakraInput,
  Collapse,
  Fade,
  FormControl,
  FormLabel,
  HStack,
  Text,
} from "@chakra-ui/react"
import { NavLink, useNavigate } from "react-router-dom"
import { bufferDecode, bufferEncode } from "../utils"
import { forwardRef, useRef, useState } from "react"

import { Props } from "../components"
import { useAuthContext } from "../App"

export const Input: React.FC<any> = forwardRef((props, ref) => {
  return <ChakraInput bgColor="white" ref={ref} {...props} />
})

export const FormField: React.FC<Props> = ({ children, ...props }) => {
  return (
    <FormControl mt="4" {...props}>
      {children}
    </FormControl>
  )
}

export type CombinedCredentials = Credential & {
  rawId: string
  response: any
}

type AuthnLoginOptions = {
  onSuccess?: () => void
  onError?: () => void
}

const handleAuthnLogin = (
  newUser: { email?: string },
  status: AuthnLoginOptions
) => {
  if (!newUser.email) return
  console.log("newUser: ", newUser)
  const formData = new FormData()
  formData.append("email", newUser.email)

  let userUid: string

  const options: RequestInit = {
    method: "POST",
    mode: "cors",
    body: formData,
    credentials: "include" /* Expecting a cookie */,
    headers: {
      Accept: "application/json",
      "Access-Control-Request-Method": "POST",
      Origin: location.origin,
      // Authorization: usernameEmail.current?.value,
    },
  }

  fetch(`http://localhost:3000/webauthn/beginlogin`, options)
    .then((res) => {
      if (res.ok) return res.json()
      throw res
    })
    .then(({ options, userUid: uid }) => {
      userUid = uid
      options.publicKey.challenge = bufferDecode(options.publicKey.challenge)
      options.publicKey.allowCredentials.forEach((listItem: any) => {
        listItem.id = bufferDecode(listItem.id)
      })

      return navigator.credentials.get({
        publicKey: options.publicKey,
      })
    })
    .then((assert) => {
      const assertion = assert as CombinedCredentials,
        authData = assertion?.response.authenticatorData,
        clientDataJSON = assertion?.response.clientDataJSON,
        rawId = assertion?.rawId,
        sig = assertion?.response.signature,
        userHandle = assertion?.response.userHandle

      fetch(`http://localhost:3000/webauthn/finishlogin/${userUid}`, {
        method: "POST",
        mode: "cors",
        body: JSON.stringify({
          id: assertion.id,
          rawId: bufferEncode(rawId as any),
          type: assertion.type,
          response: {
            authenticatorData: bufferEncode(authData),
            clientDataJSON: bufferEncode(clientDataJSON),
            signature: bufferEncode(sig),
            userHandle: bufferEncode(userHandle),
          },
        }),
        credentials: "include" /* Expecting a cookie */,
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
          "Access-Control-Request-Method": "POST",
          Origin: location.origin,
        },
      })
        .then((res) => {
          console.log("res: ", res)
          if (res.ok) {
            return res.json()
          }
          status.onError?.()
        })
        .then((user) => {
          sessionStorage.setItem("user", JSON.stringify(user))
          localStorage.setItem("user", JSON.stringify(user))
          status.onSuccess?.()
        })
        .catch(console.warn)
    })
    .catch(console.warn)
}

export function Login() {
  const navigate = useNavigate()
  const usernameEmail = useRef<HTMLInputElement>()
  const password = useRef<HTMLInputElement | undefined>()
  const account = useAuthContext()
  console.log("account: ", account)
  const [error, setError] = useState({
    invalidEmail: false,
    invalidPassword: false,
    missingEmail: false,
    missingPassword: false,
    wrongUsername: false,
    wrongPassword: false,
  })

  const handleCookieLogin:
    | React.FormEventHandler<HTMLFormElement>
    | undefined = (event) => {
    event.preventDefault()
    const formData = new FormData(event.target as HTMLFormElement)

    const options: RequestInit = {
      method: "POST",
      mode: "cors",
      body: formData,
      credentials: "include" /* Expecting a cookie */,
      headers: {
        Accept: "application/json",
        "Access-Control-Request-Method": "POST",
        Origin: location.origin,
        // Authorization: usernameEmail.current?.value,
      },
    }

    fetch(`http://localhost:3000/cookie/login`, options)
      .then((res) => {
        console.log("res: ", res)
        if (res.ok) return res.json()
        setError((p) => ({
          invalidEmail: false,
          invalidPassword: false,
          missingEmail: false,
          missingPassword: false,
          wrongUsername: true,
          wrongPassword: true,
        }))
      })
      .then((data) => {
        sessionStorage.setItem("sessionId", data.sessionId)
      })
  }

  const passwordError =
    error.invalidPassword || error.missingPassword || error.wrongPassword
  const usernameError =
    error.invalidEmail || error.missingEmail || error.wrongUsername

  return (
    <Box p="4" w="md" bg="gray.100" borderRadius="8">
      <form onSubmit={handleCookieLogin}>
        <Text fontSize="xl" fontWeight="medium">
          Login Form
        </Text>
        <FormField>
          <FormLabel>Username / Email</FormLabel>
          <Input
            ref={usernameEmail}
            defaultValue="lasse@gmail.com"
            type="text"
            name="email"
            placeholder="JohnDoe / john-doe@gmail.com.."
            autoComplete="username"
            borderColor={usernameError ? "red.600" : undefined}
            border={usernameError ? "2px solid" : undefined}
            onFocus={() => {
              setError((p) => ({
                ...p,
                invalidEmail: false,
                wrongUsername: false,
                missingEmail: false,
              }))
            }}
          />
        </FormField>
        <Fade in={usernameError}>
          <Box
            py="20px"
            px="10px"
            mt="4"
            bg="red.400"
            rounded="md"
            shadow="md"
            display={usernameError ? "block" : "none"}
          >
            {error.wrongUsername && (
              <Text>Email and username does not belong to any user</Text>
            )}
            {error.missingEmail && <Text>Field can not be empty</Text>}
            {error.invalidEmail && (
              <Text>
                The input is not a valid input, make sure your email or username
                spelled correctly
              </Text>
            )}
          </Box>
        </Fade>
        <Collapse in={false} animateOpacity>
          <FormField>
            <FormLabel>Password</FormLabel>
            <Input
              ref={password}
              type="password"
              defaultValue="secretPassword"
              name="password"
              placeholder="Secret password"
              autoComplete="current-password"
              borderColor={passwordError ? "red.600" : undefined}
              border={passwordError ? "2px solid" : undefined}
              onFocus={() => {
                setError((p) => ({
                  ...p,
                  invalidPassword: false,
                  wrongPassword: false,
                  missingPassword: false,
                }))
              }}
            />
          </FormField>
          <Fade in={passwordError}>
            <Box
              py="20px"
              px="10px"
              mt="4"
              bg="red.400"
              rounded="md"
              shadow="md"
              display={passwordError ? "block" : "none"}
            >
              {error.wrongPassword && (
                <Text>Password does not match username</Text>
              )}
              {error.missingPassword && <Text>Field can not be empty</Text>}
            </Box>
          </Fade>
        </Collapse>
        <Button colorScheme="blue" type="submit" mt="8" w="full">
          Login with Email / Password
        </Button>
        <HStack spacing="4" align="center" mt="4">
          <Button
            type="button"
            colorScheme="purple"
            w="50%"
            onClick={() => {
              handleAuthnLogin(
                { email: usernameEmail.current?.value },
                {
                  onSuccess: () => {
                    account.signIn?.()
                    navigate("/private/data")
                  },
                  onError: () => {
                    setError((p) => ({
                      invalidEmail: false,
                      invalidPassword: false,
                      missingEmail: false,
                      missingPassword: false,
                      wrongUsername: true,
                      wrongPassword: false,
                    }))
                  },
                }
              )
            }}
          >
            Login with Web Authn
          </Button>
          <Button type="button" colorScheme="orange" w="50%">
            Login with MFA
          </Button>
        </HStack>
        <NavLink to="/register">
          <Button
            colorScheme="blue"
            variant="outline"
            type="button"
            mt="4"
            w="full"
          >
            SignUp
          </Button>
        </NavLink>
      </form>
    </Box>
  )
}

// const handleUISubmit: React.FormEventHandler<HTMLFormElement> | undefined = (
//   event
// ) => {
//   event.preventDefault()
//   const formData = new FormData(event.target)
//   const data = Object.fromEntries(formData)

//   if (!data.username) setError((e) => ({ ...e, invalidEmail: true }))
//   if (!data.password) setError((e) => ({ ...e, invalidPassword: true }))

//   if (data.username && data.password)
//     setError((e) => ({ ...e, wrongPassword: true, wrongUsername: true }))
// }

// https://www.w3.org/TR/webauthn/#sctn-sample-scenarios
// const publicKey: PublicKeyCredentialCreationOptions = {
//   challenge: /* random bytes generated by the server */ Uint8Array.from(
//     "bGFzc2VBYWtqYWVyOnNlY3JldFBhc3N3b3Jk",
//     (c) => c.charCodeAt(0)
//   ),
//   rp: /* relying party */ {
//     name: "Lasse Aakjær",
//     id: "lasseaakjaer.com",
//   },
//   user: {
//     id: Uint8Array.from(
//       window.atob("MIIBkzCCATigAwIBAjCCAZMwggE4oAMCAQIwggGTMII="),
//       (c) => c.charCodeAt(0)
//     ),
//     name: "alex.mueller@example.com",
//     displayName: "Alex Müller",
//   },
//   pubKeyCredParams: [
//     {
//       type: "public-key",
//       alg: -7, // "ES256" as registered in the IANA COSE Algorithms registry
//     },
//     {
//       type: "public-key",
//       alg: -257, // Value registered by this specification for "RS256"
//     },
//   ],
//   authenticatorSelection: {
//     /*  https://w3c.github.io/webauthn/#dom-authenticatorselectioncriteria-userverification
//       This enumeration’s values describe authenticators' attachment modalities.
//       Relying Parties use this to express a preferred authenticator attachment modality when calling navigator.credentials.create() to create a credential,
//       and clients use this to report the authenticator attachment modality used to complete a registration or authentication ceremony.
//     */
//     // Try to use UV (User verification) if possible. This is also the default.
//     userVerification: "preferred" /*   "cross-platform" */,
//     /*   Authenticator Attachment Modality
//     A platform authenticator is attached using a client device-specific transport, called platform attachment, and is usually not removable from the client device.
//     A public key credential bound to a platform authenticator is called a platform credential.

//     A roaming authenticator is attached using cross-platform transports, called cross-platform attachment.
//     Authenticators of this class are removable from, and can "roam" between, client devices.
//     A public key credential bound to a roaming authenticator is called a roaming credential.
//     */
//   },
//   timeout: 360000, // 6 minutes
//   excludeCredentials: [
//     // Don’t re-register any authenticator that has one of these credentials
//     {
//       id: Uint8Array.from(
//         window.atob("ufJWp8YGlibm1Kd9XQBWN1WAw2jy5In2Xhon9HAqcXE="),
//         (c) => c.charCodeAt(0)
//       ),
//       type: "public-key",
//     },
//     {
//       id: Uint8Array.from(
//         window.atob("E/e1dhZc++mIsz4f9hb6NifAzJpF1V4mEtRlIPBiWdY="),
//         (c) => c.charCodeAt(0)
//       ),
//       type: "public-key",
//     },
//   ],

//   // Make excludeCredentials check backwards compatible with credentials registered with U2F
//   extensions: { appidExclude: "https://acme.example.com" },
// }
