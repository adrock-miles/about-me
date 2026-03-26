import { Box, Flex, Heading, Text, Separator, Badge } from "@radix-ui/themes";
import { useProjects } from "../hooks/useProjects";

function Projects() {
  const { projects, loading, error } = useProjects();

  return (
    <Box className="page fade-in">
      <Flex direction="column" gap="6" py="8">
        <Box>
          <Text size="2" color="gray" mb="2" style={{ display: "block" }}>
            Recent Work
          </Text>
          <Heading
            size="7"
            weight="medium"
            className="page-title"
            style={{ letterSpacing: "-0.03em" }}
          >
            Projects
          </Heading>
          <Text size="3" color="gray" style={{ maxWidth: 580, lineHeight: 1.7 }}>
            A selection of projects I've recently built or contributed to. Each
            one represents a problem I wanted to solve or a technology I wanted
            to explore.
          </Text>
        </Box>

        <Separator size="4" />

        {loading && (
          <Text size="2" color="gray">
            Loading projects...
          </Text>
        )}

        {error && (
          <Text size="2" color="red">
            {error}
          </Text>
        )}

        <Flex direction="column" gap="5">
          {projects.map((project) => (
            <Box
              key={project.name}
              p="5"
              style={{
                border: "1px solid var(--gray-a4)",
                borderRadius: "var(--radius-3)",
                transition: "border-color 0.2s",
              }}
            >
              <Flex direction="column" gap="3">
                <Flex justify="between" align="start">
                  <Flex direction="column" gap="1">
                    <Text size="3" weight="medium" asChild>
                      <a
                        href={project.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        style={{ borderBottom: "1px solid var(--gray-a6)" }}
                      >
                        {project.name}
                      </a>
                    </Text>
                    {project.description && (
                      <Text size="2" color="gray" style={{ lineHeight: 1.6 }}>
                        {project.description}
                      </Text>
                    )}
                  </Flex>
                  {project.stars > 0 && (
                    <Text size="1" color="gray">
                      {project.stars} stars
                    </Text>
                  )}
                </Flex>

                <Flex gap="2" wrap="wrap">
                  {project.language && (
                    <Badge variant="soft" color="gray" size="1">
                      {project.language}
                    </Badge>
                  )}
                  {project.topics?.map((topic) => (
                    <Badge key={topic} variant="outline" color="gray" size="1">
                      {topic}
                    </Badge>
                  ))}
                </Flex>
              </Flex>
            </Box>
          ))}
        </Flex>
      </Flex>
    </Box>
  );
}

export default Projects;
