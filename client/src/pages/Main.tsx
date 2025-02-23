import { Button } from "@/components/ui/button.tsx";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx";
import { Input } from "@/components/ui/input.tsx";
import { FormEventHandler } from "react";
import {
  InstanceSchema,
  useAddInstance,
  useInstances,
  useKillInstance,
  useModifyInstance,
} from "./hooks.tsx";
import { Ban, Box, RadioReceiver, SatelliteDish } from "lucide-react";
import { toast } from "sonner";

interface AddInstanceForm extends HTMLFormElement {
  instance_name: HTMLInputElement;
}

function AddInstance() {
  const { mutate } = useAddInstance();
  const handleSubmit: FormEventHandler<AddInstanceForm> = (e) => {
    e.preventDefault();
    const name = e.currentTarget.instance_name.value;
    mutate({
      name,
    });
  };
  return (
    <>
      <Dialog>
        <DialogTrigger asChild>
          <Button className="mb-8">Add Instance</Button>
        </DialogTrigger>
        <DialogContent>
          <form onSubmit={handleSubmit}>
            <DialogHeader>
              <DialogTitle>Add Instance</DialogTitle>
              <DialogDescription>
                <Input type="text" id="instance_name"></Input>
              </DialogDescription>
              <DialogFooter>
                <Button type="submit">Add</Button>
              </DialogFooter>
            </DialogHeader>
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}

function Instance(props: { data: InstanceSchema; standby?: string }) {
  const { data, standby } = props;
  const { mutate: kill } = useKillInstance(data.Name);
  const { mutate: modify } = useModifyInstance(data.Name);

  return (
    <Card className="my-4">
      <CardHeader>
        <CardTitle>{data.Name}</CardTitle>
        <CardDescription>Mapped to port {data.Port}</CardDescription>
      </CardHeader>
      <CardContent>
        <Button
          className="ml-2"
          onClick={() =>
            modify({
              Primary: true,
            })
          }
        >
          <SatelliteDish />
          Setup as Primary
        </Button>
        <Button
          className="mx-2"
          onClick={() => {
            if (standby == undefined || standby == "") {
              toast.error("Cannot set as standby without a leader!");
              return;
            }

            modify({
              Standby: true,
              StandbyTo: standby,
            });
          }}
        >
          <RadioReceiver />
          Setup as Standby
        </Button>
        <Button className="mr-2" onClick={() => kill()}>
          <Ban />
          Kill
        </Button>
      </CardContent>
    </Card>
  );
}

function Main() {
  const { data } = useInstances();
  if (data == undefined) {
    return null;
  }
  return (
    <div className="p-16">
      <AddInstance />
      {data.map((x) => (
        <Instance key={x.Name} data={x} />
      ))}
    </div>
  );
}
export default Main;
