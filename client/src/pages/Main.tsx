import { Button } from "@/components/ui/button.tsx";
import { useInstances } from "./hooks.tsx";

function Main() {
  const { data } = useInstances();
  if (data == undefined) {
    return null;
  }
  return (
    <div>
      <Button>Yo</Button>
      {data.map((x) => (
        <div>{x.Name}</div>
      ))}
    </div>
  );
}
export default Main;
