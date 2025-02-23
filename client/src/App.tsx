import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "./App.css";
import Main from "./pages/Main.tsx";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Main />
    </QueryClientProvider>
  );
}

export default App;
