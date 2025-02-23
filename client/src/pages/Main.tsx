import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTrigger,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog.tsx";
import { useAddInstance, useInstances } from "./hooks.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Input } from "@/components/ui/input.tsx";
import { FormEventHandler } from "react";

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
        <DialogTrigger>Add Instance</DialogTrigger>
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

function Main() {
  const { data } = useInstances();
  if (data == undefined) {
    return null;
  }
  return (
    <div>
      <AddInstance />
      {data.map((x) => (
        <div>{x.Name}</div>
      ))}
    </div>
  );
}
export default Main;
