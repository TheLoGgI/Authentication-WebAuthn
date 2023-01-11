import { Box } from "@chakra-ui/react"
import React from "react"

export type Props = {
  children?: React.ReactNode
}

export const Layout: React.FC<Props> = ({ children }) => {
  return (
    <Box
      as="main"
      display="flex"
      w="full"
      h="90vh"
      justifyContent="center"
      alignItems="center"
    >
      {children}
    </Box>
  )
}
