import queryClient from "@/client.ts";
import { useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { z } from "zod";

const instanceSchema = z.object({
  Name: z.string(),
  ContainerID: z.string(),
  NetworkID: z.string(),
  Port: z.string(),
});

export type InstanceSchema = z.infer<typeof instanceSchema>;

export function useInstances() {
  return useQuery({
    queryKey: ["instances"],
    queryFn: async () => {
      const res = await fetch("/api/instances");
      const data = await res.json();
      return instanceSchema.array().parse(data);
    },
  });
}

export function useAddInstance() {
  return useMutation({
    mutationFn: async (body: { name: string }) => {
      const res = await fetch("/api/instances", {
        method: "post",
        body: JSON.stringify(body),
      });
      const data = await res.json();
      if (res.status !== 200) {
        throw new Error(data.Message);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["instances"],
      });
      toast.success("Instance created!");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}

export function useKillInstance(name: string) {
  return useMutation({
    mutationFn: async () => {
      const res = await fetch(`/api/instances/${name}`, {
        method: "delete",
      });
      const data = await res.json();
      if (res.status !== 200) {
        throw new Error(data.Message);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["instances"],
      });
      toast.success("Instance deleted!");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}

type ModifyBody =
  | {
      Primary: boolean;
    }
  | {
      Standby: boolean;
      StandbyTo: string;
    };

export function useModifyInstance(name: string) {
  return useMutation({
    mutationFn: async (body: ModifyBody) => {
      const res = await fetch(`/api/instances/${name}`, {
        method: "put",
        body: JSON.stringify(body),
      });
      const data = await res.json();
      if (res.status !== 200) {
        throw new Error(data.Message);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["instances"],
      });
      toast.success("Instance modified!");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}
