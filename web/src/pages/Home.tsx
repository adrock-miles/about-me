import { Box, Flex, Heading, Text, Separator } from "@radix-ui/themes";
import { Link } from "react-router-dom";

function Home() {
  return (
    <Box className="page fade-in">
      <Flex direction="column" gap="6" py="8">
        <Box>
          <Text size="2" color="gray" mb="2" style={{ display: "block" }}>
            Software Engineer
          </Text>
          <Heading size="8" weight="medium" style={{ letterSpacing: "-0.03em" }}>
            Hi, I'm Adam Miles.
          </Heading>
        </Box>

        <Text size="3" color="gray" style={{ maxWidth: 580, lineHeight: 1.7 }}>
          I'm a software engineer who builds thoughtful, well-crafted systems.
          I care about clean architecture, developer experience, and writing
          code that's easy to reason about. I'm always exploring new
          technologies and pushing myself to grow.
        </Text>

        <Separator size="4" />

        <Flex direction="column" gap="4">
          <Heading size="4" weight="medium">
            What I Do
          </Heading>
          <Flex gap="8" wrap="wrap">
            <Flex direction="column" gap="1" style={{ maxWidth: 260 }}>
              <Text size="2" weight="medium">
                Backend Systems
              </Text>
              <Text size="2" color="gray">
                Designing scalable APIs and services with Go, following
                domain-driven design principles.
              </Text>
            </Flex>
            <Flex direction="column" gap="1" style={{ maxWidth: 260 }}>
              <Text size="2" weight="medium">
                Frontend Development
              </Text>
              <Text size="2" color="gray">
                Building modern, accessible interfaces with React and
                TypeScript.
              </Text>
            </Flex>
            <Flex direction="column" gap="1" style={{ maxWidth: 260 }}>
              <Text size="2" weight="medium">
                DevOps & Infrastructure
              </Text>
              <Text size="2" color="gray">
                Containerization, CI/CD pipelines, and cloud infrastructure
                management.
              </Text>
            </Flex>
          </Flex>
        </Flex>

        <Separator size="4" />

        <Flex gap="4">
          <Text size="2" weight="medium" asChild>
            <Link to="/projects" style={{ borderBottom: "1px solid var(--gray-a8)" }}>
              View my projects
            </Link>
          </Text>
          <Text size="2" weight="medium" asChild>
            <Link to="/contact" style={{ borderBottom: "1px solid var(--gray-a8)" }}>
              Get in touch
            </Link>
          </Text>
        </Flex>
      </Flex>
    </Box>
  );
}

export default Home;
