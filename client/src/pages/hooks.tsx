import { useQuery } from "@tanstack/react-query";
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
    queryKey: ["/instances"],
    queryFn: async () => {
      const res = await fetch("/api/instances");
      const data = await res.json();
      return InstanceSchema.parse(data);
    },
  });
}
