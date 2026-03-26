import { useState } from "react";

interface ContactForm {
  name: string;
  email: string;
  subject: string;
  body: string;
}

export function useContact() {
  const [sending, setSending] = useState(false);
  const [sent, setSent] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const send = async (form: ContactForm) => {
    setSending(true);
    setError(null);

    try {
      const res = await fetch("/api/contact", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(form),
      });

      if (!res.ok) throw new Error("Failed to send message");
      setSent(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setSending(false);
    }
  };

  return { send, sending, sent, error, reset: () => setSent(false) };
}
