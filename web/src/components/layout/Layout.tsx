import { ReactNode } from "react";
import { Flex } from "@radix-ui/themes";
import Nav from "./Nav";
import Footer from "./Footer";

interface LayoutProps {
  children: ReactNode;
}

function Layout({ children }: LayoutProps) {
  return (
    <Flex direction="column" minHeight="100vh">
      <Nav />
      <Flex flexGrow="1" direction="column">
        {children}
      </Flex>
      <Footer />
    </Flex>
  );
}

export default Layout;
