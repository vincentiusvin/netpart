import { QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "sonner";
import queryClient from "./client.ts";
import "./index.css";
import Main from "./pages/Main.tsx";

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Main />
      <Toaster />
    </QueryClientProvider>
  );
}

export default App;
