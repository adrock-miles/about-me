import { Flex, Text, Separator } from "@radix-ui/themes";

function Footer() {
  return (
    <Flex direction="column" py="6">
      <Separator size="4" mb="4" />
      <Flex justify="between" align="center">
        <Text size="1" color="gray">
          &copy; {new Date().getFullYear()} Miles Adrock
        </Text>
        <Flex gap="4">
          <Text size="1" color="gray" asChild>
            <a
              href="https://github.com/adrock-miles"
              target="_blank"
              rel="noopener noreferrer"
            >
              GitHub
            </a>
          </Text>
          <Text size="1" color="gray" asChild>
            <a href="mailto:miles.adrock@gmail.com">Email</a>
          </Text>
        </Flex>
      </Flex>
    </Flex>
  );
}

export default Footer;
