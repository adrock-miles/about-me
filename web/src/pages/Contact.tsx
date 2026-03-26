import { useState, FormEvent } from "react";
import {
  Box,
  Flex,
  Heading,
  Text,
  Separator,
  TextField,
  TextArea,
  Button,
} from "@radix-ui/themes";
import { useContact } from "../hooks/useContact";

function Contact() {
  const { send, sending, sent, error } = useContact();
  const [form, setForm] = useState({
    name: "",
    email: "",
    subject: "",
    body: "",
  });

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    send(form);
  };

  const update = (field: string, value: string) => {
    setForm((prev) => ({ ...prev, [field]: value }));
  };

  if (sent) {
    return (
      <Box className="page fade-in">
        <Flex direction="column" gap="4" py="8" align="center">
          <Heading size="6" weight="medium">
            Message sent.
          </Heading>
          <Text size="3" color="gray">
            Thanks for reaching out. I'll get back to you soon.
          </Text>
        </Flex>
      </Box>
    );
  }

  return (
    <Box className="page fade-in">
      <Flex direction="column" gap="6" py="8">
        <Box>
          <Text size="2" color="gray" mb="2" style={{ display: "block" }}>
            Let's Connect
          </Text>
          <Heading
            size="7"
            weight="medium"
            className="page-title"
            style={{ letterSpacing: "-0.03em" }}
          >
            Contact
          </Heading>
          <Text size="3" color="gray" style={{ maxWidth: 580, lineHeight: 1.7 }}>
            Have a question, opportunity, or just want to say hello? Drop me a
            message and I'll get back to you as soon as I can.
          </Text>
        </Box>

        <Separator size="4" />

        <Box style={{ maxWidth: 480 }}>
          <form onSubmit={handleSubmit}>
            <Flex direction="column" gap="4">
              <Flex direction="column" gap="1">
                <Text as="label" size="2" weight="medium" htmlFor="name">
                  Name
                </Text>
                <TextField.Root
                  id="name"
                  placeholder="Your name"
                  value={form.name}
                  onChange={(e) => update("name", e.target.value)}
                  required
                />
              </Flex>

              <Flex direction="column" gap="1">
                <Text as="label" size="2" weight="medium" htmlFor="email">
                  Email
                </Text>
                <TextField.Root
                  id="email"
                  type="email"
                  placeholder="you@example.com"
                  value={form.email}
                  onChange={(e) => update("email", e.target.value)}
                  required
                />
              </Flex>

              <Flex direction="column" gap="1">
                <Text as="label" size="2" weight="medium" htmlFor="subject">
                  Subject
                </Text>
                <TextField.Root
                  id="subject"
                  placeholder="What's this about?"
                  value={form.subject}
                  onChange={(e) => update("subject", e.target.value)}
                />
              </Flex>

              <Flex direction="column" gap="1">
                <Text as="label" size="2" weight="medium" htmlFor="body">
                  Message
                </Text>
                <TextArea
                  id="body"
                  placeholder="Your message..."
                  value={form.body}
                  onChange={(e) => update("body", e.target.value)}
                  rows={6}
                  required
                />
              </Flex>

              {error && (
                <Text size="2" color="red">
                  {error}
                </Text>
              )}

              <Button
                type="submit"
                variant="solid"
                color="gray"
                highContrast
                disabled={sending}
                style={{ cursor: "pointer" }}
              >
                {sending ? "Sending..." : "Send Message"}
              </Button>
            </Flex>
          </form>
        </Box>
      </Flex>
    </Box>
  );
}

export default Contact;
