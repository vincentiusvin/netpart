import { QueryClientProvider } from "@tanstack/react-query";
import queryClient from "./client.ts";
import "./index.css";
import Main from "./pages/Main.tsx";
import { ThemeProvider } from "./pages/theme-provider.tsx";
import { Toaster } from "sonner";

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider defaultTheme="dark">
        <Main />
        <Toaster />
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App;
