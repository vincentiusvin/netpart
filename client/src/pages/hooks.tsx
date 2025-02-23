import queryClient from "@/client.ts";
import { QueryClient, useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { z } from "zod";

const InstanceSchema = z
  .object({
    Name: z.string(),
    ContainerID: z.string(),
    NetworkID: z.string(),
    Port: z.string(),
  })
  .array();

export function useInstances() {
  return useQuery({
    queryKey: ["instances"],
    queryFn: async () => {
      const res = await fetch("/api/instances");
      const data = await res.json();
      return InstanceSchema.parse(data);
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
