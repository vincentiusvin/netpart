import { Button } from "@/components/ui/button.tsx";
import {
  Card,
  CardContent,
  CardDescription,
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
import { Label } from "@/components/ui/label.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx";
import {
  Ban,
  BetweenHorizonalEnd,
  Plus,
  RadioReceiver,
  SatelliteDish,
} from "lucide-react";
import { FormEventHandler } from "react";
import { toast } from "sonner";
import {
  InstanceSchema,
  useAddInstance,
  useInstanceData,
  useInstances,
  useKillInstance,
  useModifyInstance,
  usePutInstanceData,
} from "./hooks.tsx";

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
          <Button className="mb-8">
            <Plus />
            Add Instance
          </Button>
        </DialogTrigger>
        <DialogContent>
          <form onSubmit={handleSubmit}>
            <DialogHeader>
              <DialogTitle>Add Instance</DialogTitle>
              <DialogDescription className="my-2">
                <Label htmlFor="text">Name</Label>
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
          className="mr-2"
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
        <Button className="mx-2" onClick={() => kill()}>
          <Ban />
          Kill
        </Button>
        <div>
          <DataSubmission data={data} />
        </div>
        <Data data={data} />
      </CardContent>
    </Card>
  );
}

interface DataSubmissionForm extends HTMLFormElement {
  datakey: HTMLInputElement;
  datavalue: HTMLInputElement;
}

function DataSubmission(props: { data: InstanceSchema }) {
  const { data } = props;
  const { mutate: put } = usePutInstanceData(data.Name);

  const handleSubmit: FormEventHandler<DataSubmissionForm> = (e) => {
    e.preventDefault();
    const key = e.currentTarget.datakey.value;
    const value = e.currentTarget.datavalue.value;
    put({
      Key: key,
      Value: value,
    });
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="flex gap-8 items-center">
        <div>
          <Label>Key</Label>
          <Input type="text" id="datakey" />
          <Label>Value</Label>
          <Input type="text" id="datavalue" />
        </div>
        <Button className="mr-2 mt-4" type="submit">
          <BetweenHorizonalEnd />
          Submit
        </Button>
      </div>
    </form>
  );
}

function Data(props: { data: InstanceSchema }) {
  const { data } = props;
  const { data: kvs } = useInstanceData(data.Name);

  if (kvs == undefined) {
    return (
      <Table className="mt-4">
        <TableCaption>Data for {data.Name}</TableCaption>
        <Skeleton />
      </Table>
    );
  }

  return (
    <Table className="mt-4">
      <TableHeader>
        <TableRow>
          <TableHead>Key</TableHead>
          <TableHead>Value</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {kvs.map((x) => (
          <TableRow key={x.Key}>
            <TableCell>{x.Key}</TableCell>
            <TableCell>{x.Value}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function Main() {
  const { data } = useInstances();
  if (data == undefined) {
    return <Skeleton />;
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
