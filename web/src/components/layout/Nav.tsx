import { Link, useLocation } from "react-router-dom";
import { Flex, Text } from "@radix-ui/themes";

const navItems = [
  { path: "/", label: "Home" },
  { path: "/interests", label: "Interests" },
  { path: "/projects", label: "Projects" },
  { path: "/contact", label: "Contact" },
];

function Nav() {
  const location = useLocation();

  return (
    <Flex
      asChild
      justify="between"
      align="center"
      py="4"
      style={{ borderBottom: "1px solid var(--gray-a4)" }}
    >
      <nav>
        <Text size="4" weight="medium" style={{ letterSpacing: "-0.02em" }}>
          <Link to="/">miles adrock</Link>
        </Text>
        <Flex gap="6">
          {navItems.map((item) => (
            <Text
              key={item.path}
              size="2"
              weight={location.pathname === item.path ? "medium" : "regular"}
              style={{
                opacity: location.pathname === item.path ? 1 : 0.6,
                transition: "opacity 0.2s",
              }}
            >
              <Link to={item.path}>{item.label}</Link>
            </Text>
          ))}
        </Flex>
      </nav>
    </Flex>
  );
}

export default Nav;
