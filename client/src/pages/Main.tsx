import { Button } from "@/components/ui/button.tsx";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox.tsx";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";
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
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs.tsx";
import { TabsContent } from "@radix-ui/react-tabs";
import {
  Ban,
  BetweenHorizonalEnd,
  Plus,
  RadioReceiver,
  SatelliteDish,
} from "lucide-react";
import React, { FormEventHandler, useState } from "react";
import {
  InstanceSchema,
  useAddInstance,
  useConnect,
  useDisconnect,
  useGetConnection,
  useInstanceData,
  useInstanceReplication,
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
  );
}

function StandbyDialog(props: { data: InstanceSchema }) {
  const { data } = props;
  const { mutate: modify } = useModifyInstance(data.Name);
  const { data: leaderOptions } = useInstances();
  const [leader, setLeader] = useState("");
  const [enable, setEnable] = useState(false);

  const handleSubmit: FormEventHandler = (e) => {
    e.preventDefault();
    if (enable) {
      modify({
        Refresh: true,
        RefreshTo: leader,
      });
    } else {
      modify({
        Standby: true,
        StandbyTo: leader,
      });
    }
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button className="mx-2">
          <RadioReceiver />
          Setup as Standby
        </Button>
      </DialogTrigger>
      <DialogContent>
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Setup as Standby</DialogTitle>
            <DialogDescription className="my-2">
              <Label htmlFor="text">Leader</Label>
              <Select onValueChange={(v) => setLeader(v)}>
                <SelectTrigger>
                  <SelectValue placeholder="Select a leader" />
                </SelectTrigger>
                <SelectContent>
                  {leaderOptions?.map((x) => (
                    <SelectItem value={x.Name} key={x.Name}>
                      {x.Name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <br />
              <Label htmlFor="text">Refresh: </Label>
              <Checkbox
                checked={enable}
                onClick={() => {
                  setEnable((x) => !x);
                }}
              />
            </DialogDescription>
            <DialogFooter>
              <Button type="submit">Submit</Button>
            </DialogFooter>
          </DialogHeader>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function Instance(props: { data: InstanceSchema }) {
  const { data } = props;
  const { mutate: kill } = useKillInstance(data.Name);
  const { mutate: modify } = useModifyInstance(data.Name);
  const { data: repl } = useInstanceReplication(data.Name, 1000);

  const actives =
    repl?.ActiveData.map((x) => ({
      icon: <SatelliteDish />,
      title: "Primary for " + x.Application_Name,
      subtitle: x.State + " - " + x.Sync_State,
    })) || [];

  const passives =
    repl?.StandbyData.map((x) => ({
      icon: <RadioReceiver />,
      title: "Standby for " + x.Subname,
      subtitle: x.Subenabled ? "Enabled" : "Disabled",
    })) || [];

  const replData: {
    icon: React.ReactNode;
    title: string;
    subtitle: string;
  }[] = actives.concat(passives);

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
        <StandbyDialog data={data} />
        <Button className="mx-2" onClick={() => kill()}>
          <Ban />
          Kill
        </Button>
        <div className="mt-4">
          {replData.map((x) => (
            <div key={x.title} className="flex gap-4 items-center mt-2">
              {x.icon}
              <div>
                <div className="font-semibold">{x.title}</div>
                <div className="text-muted-foreground">{x.subtitle}</div>
              </div>
            </div>
          ))}
        </div>
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

function DBManager(props: { data: InstanceSchema }) {
  const { data } = props;
  const { data: kvs } = useInstanceData(data.Name, 500);

  if (kvs == undefined) {
    return (
      <Table className="mt-4">
        <TableCaption>Data for {data.Name}</TableCaption>
      </Table>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{data.Name}</CardTitle>
      </CardHeader>
      <CardContent>
        <DataSubmission data={data} />
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
      </CardContent>
    </Card>
  );
}

function Data() {
  const { data: instances } = useInstances();

  if (instances == undefined) {
    return <Skeleton />;
  }

  return (
    <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-4">
      {instances.map((x) => (
        <DBManager data={x} />
      ))}
    </div>
  );
}

function NetworkToggle(props: { instance1: string; instance2: string }) {
  const { instance1, instance2 } = props;
  const { data: connect } = useGetConnection(instance1, instance2);
  const { mutate: doConnect } = useConnect(instance1, instance2);
  const { mutate: doDisconnect } = useDisconnect(instance1, instance2);
  if (connect == undefined) {
    return <Skeleton />;
  }

  return (
    <Checkbox
      checked={connect.Connected}
      onClick={() => {
        if (connect.Connected) {
          doDisconnect();
        } else {
          doConnect();
        }
      }}
    />
  );
}

function Network() {
  const { data } = useInstances();
  if (data == undefined) {
    return <Skeleton />;
  }

  const matrix = data.map((x, i) => {
    return data.map((y, j) => ({
      Enabled: i < j,
      Node1: x.Name,
      Node2: y.Name,
    }));
  });

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead />
          {data.map((x) => (
            <TableHead key={x.Name} className="text-center align-middle">
              {x.Name}
            </TableHead>
          ))}
        </TableRow>
      </TableHeader>
      <TableBody>
        {matrix.map((x, i) => (
          <TableRow key={i}>
            <TableCell
              className={
                "text-muted-foreground h-10 px-2 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]"
              }
            >
              {data[i].Name}
            </TableCell>
            {x.map((y, j) => (
              <TableCell className="text-center" key={j}>
                {y.Enabled ? (
                  <NetworkToggle instance1={y.Node1} instance2={y.Node2} />
                ) : (
                  "Redundant"
                )}
              </TableCell>
            ))}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function Provision() {
  const { data } = useInstances();
  if (data == undefined) {
    return <Skeleton />;
  }
  return (
    <>
      <AddInstance />
      {data.map((x) => (
        <Instance key={x.Name} data={x} />
      ))}
    </>
  );
}

function Main() {
  return (
    <Tabs defaultValue="provision" className="p-16">
      <TabsList className="mb-8">
        <TabsTrigger value="provision">Provision</TabsTrigger>
        <TabsTrigger value="network">Network</TabsTrigger>
        <TabsTrigger value="data">Data</TabsTrigger>
      </TabsList>
      <TabsContent value="provision">
        <Provision />
      </TabsContent>
      <TabsContent value="network">
        <Network />
      </TabsContent>
      <TabsContent value="data">
        <Data />
      </TabsContent>
    </Tabs>
  );
  return <Network />;
}
export default Main;
