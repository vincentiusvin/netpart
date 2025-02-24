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

const dataSchema = z.object({
  Key: z.string(),
  Value: z.string(),
});

export function useInstanceData(name: string, refetch?: number) {
  return useQuery({
    queryKey: ["instances", name, "keys"],
    refetchInterval: refetch,
    queryFn: async () => {
      const res = await fetch(`/api/instances/${name}/keys`);
      const data = await res.json();
      return dataSchema.array().parse(data);
    },
  });
}

const replicationSchema = z.object({
  ActiveData: z
    .object({
      Application_Name: z.string(),
      State: z.string(),
      Sync_State: z.string(),
    })
    .array(),
  StandbyData: z
    .object({
      Subname: z.string(),
      Subenabled: z.boolean(),
    })
    .array(),
});

export function useInstanceReplication(name: string, refetch?: number) {
  return useQuery({
    queryKey: ["instances", name, "replication"],
    refetchInterval: refetch,
    queryFn: async () => {
      const res = await fetch(`/api/instances/${name}`);
      const data = await res.json();
      return replicationSchema.parse(data);
    },
  });
}

type PutBody = {
  Key: string;
  Value: string;
};

export function usePutInstanceData(name: string) {
  return useMutation({
    mutationFn: async (body: PutBody) => {
      const res = await fetch(`/api/instances/${name}/keys/${body.Key}`, {
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
        queryKey: ["instances", name, "keys"],
      });
      toast.success("Instance modified!");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}

const networkSchema = z.object({
  Connected: z.boolean(),
  Message: z.string(),
});

export function useGetConnection(name1: string, name2: string) {
  return useQuery({
    queryKey: ["instances", name1, "net", name2],
    queryFn: async () => {
      const res = await fetch(`/api/instances/${name1}/connections/${name2}`);
      const data = await res.json();
      return networkSchema.parse(data);
    },
  });
}

export function useConnect(name1: string, name2: string) {
  return useMutation({
    mutationFn: async () => {
      const res = await fetch(`/api/instances/${name1}/connections/${name2}`, {
        method: "put",
      });
      const data = await res.json();
      if (res.status !== 200) {
        throw new Error(data.Message);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["instances", name1, "net", name2],
      });
      toast.success("Network connected!");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}

export function useDisconnect(name1: string, name2: string) {
  return useMutation({
    mutationFn: async () => {
      const res = await fetch(`/api/instances/${name1}/connections/${name2}`, {
        method: "delete",
      });
      const data = await res.json();
      if (res.status !== 200) {
        throw new Error(data.Message);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["instances", name1, "net", name2],
      });
      toast.success("Network disconnected!");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}
