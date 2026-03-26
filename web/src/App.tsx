import { Routes, Route } from "react-router-dom";
import { Container } from "@radix-ui/themes";
import Layout from "./components/layout/Layout";
import Home from "./pages/Home";
import Interests from "./pages/Interests";
import Projects from "./pages/Projects";
import Contact from "./pages/Contact";

function App() {
  return (
    <Container size="3" px="4">
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/interests" element={<Interests />} />
          <Route path="/projects" element={<Projects />} />
          <Route path="/contact" element={<Contact />} />
        </Routes>
      </Layout>
    </Container>
  );
}

export default App;
