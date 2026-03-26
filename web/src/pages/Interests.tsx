import { Box, Flex, Heading, Text, Separator } from "@radix-ui/themes";

interface InterestSection {
  title: string;
  description: string;
}

const interests: InterestSection[] = [
  {
    title: "Open Source",
    description:
      "I enjoy contributing to and building open source projects. It's a great way to collaborate with developers worldwide and give back to the community that has given me so much.",
  },
  {
    title: "Systems Design",
    description:
      "I'm fascinated by distributed systems, clean architecture patterns, and the challenge of building software that scales gracefully. Domain-driven design is a core part of how I think about problems.",
  },
  {
    title: "Developer Tooling",
    description:
      "I love building tools that make developers' lives easier. Whether it's CLI utilities, automation scripts, or development environment setups, I believe great tooling amplifies productivity.",
  },
  {
    title: "Continuous Learning",
    description:
      "Technology moves fast, and I enjoy keeping up. I'm always diving into new languages, frameworks, and paradigms to broaden my perspective and sharpen my skills.",
  },
  {
    title: "Music & Creativity",
    description:
      "Outside of code, I find inspiration in music and creative pursuits. I believe creativity in engineering and the arts are deeply connected — both require craft, iteration, and a willingness to experiment.",
  },
  {
    title: "Fitness & Wellness",
    description:
      "I believe in maintaining a healthy balance. Regular exercise and mindfulness help me stay focused and bring my best energy to everything I do.",
  },
];

function Interests() {
  return (
    <Box className="page fade-in">
      <Flex direction="column" gap="6" py="8">
        <Box>
          <Text size="2" color="gray" mb="2" style={{ display: "block" }}>
            Beyond the Code
          </Text>
          <Heading
            size="7"
            weight="medium"
            className="page-title"
            style={{ letterSpacing: "-0.03em" }}
          >
            Interests & Passions
          </Heading>
          <Text size="3" color="gray" style={{ maxWidth: 580, lineHeight: 1.7 }}>
            A look at what drives me, both professionally and personally.
          </Text>
        </Box>

        <Separator size="4" />

        <Flex direction="column" gap="6">
          {interests.map((interest, i) => (
            <Flex key={i} direction="column" gap="2">
              <Text size="3" weight="medium">
                {interest.title}
              </Text>
              <Text
                size="2"
                color="gray"
                style={{ maxWidth: 580, lineHeight: 1.7 }}
              >
                {interest.description}
              </Text>
            </Flex>
          ))}
        </Flex>
      </Flex>
    </Box>
  );
}

export default Interests;
