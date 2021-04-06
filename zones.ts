import {
  json,
  serve,
  validateRequest,
} from "https://deno.land/x/sift@0.1.7/mod.ts";
import { config } from "https://deno.land/x/dotenv@v2.0.0/mod.ts";

interface Temperature {
  value: number
  unit: 'CELSIUS' | 'FAHRENHEIT'
}

interface Zone {
  name: string
  temperature: Temperature
  humidity: number
}

serve({
  "/zones": handleZones,
});


async function handleZones(request: Request) {
  // We allow GET requests and POST requests with the
  // following fields ("name", "temperature", "humidity") in the body.
  const { error, body } = await validateRequest(request, {
    GET: {},
    POST: {
      body: ["name"],
    },
  });
  // validateRequest populates the error if the request doesn't meet
  // the schema we defined.
  if (error) {
    return json({ error: error.message }, { status: error.status });
  }

  // Handle POST requests.
  if (request.method === "POST") {
    const { name, temperature, humidity, error } = await createZone(
      (body as {name: string}).name,
    );
    if (error) {
      return json({ error: "couldn't create the zone" }, { status: 500 });
    }

    return json({ name, temperature, humidity }, { status: 201 });
  }

  // Handle GET requests.
  {
    const { zones, error } = await getAllZones();
    if (error) {
      return json({ error: "couldn't fetch the zones" }, { status: 500 });
    }

    return json({ zones });
  }
}

/** Get all zones available in the database. */
async function getAllZones() {
  const query = `
    query {
      allZones {
        data {
          name
          temperature
          humidity
        }
      }
    }
  `;

  const {
    data: {
      allZones: { data: zones },
    },
    error,
  } = await queryFuana(query, {});
  if (error) {
    return { error };
  }

  return { zones };
}

/** Create a new zone in the database. */
async function createZone(name: string): Promise<{ name?: string, temperature?: Temperature; humidity?: number; error?: string }> {
  const query = `
    mutation($name: String!) {
      createZone(
        data: { name: $name, temperature: { create: { value: 0, unit: CELCIUS } }, humidity: 0}
      ) {
        _id
        name
        temperature {
            value
            unit
        }
        humidity
      }
    }
  `;
    
  const {
    data: { createZone },
    error,
  } = await queryFuana(query, { name });
  if (error) {
    return { error };
  }

  return createZone; // Zone
}

async function queryFuana(
  query: string,
  variables: { [key: string]: unknown },
): Promise<{
  data?: any;
  error?: any;
}> {
  // Grab the secret from the environment.
  const token = config().FAUNA_SECRET;
  if (!token) {
    throw new Error("environment variable FAUNA_SECRET not set");
  }

  try {
    // Make a POST request to fauna's graphql endpoint with body being
    // the query and its variables.
    const res = await fetch("https://graphql.fauna.com/graphql", {
      method: "POST",
      headers: {
        authorization: `Bearer ${token}`,
        "content-type": "application/json",
      },
      body: JSON.stringify({
        query,
        variables,
      }),
    });

    const { data, errors } = await res.json();
    if (errors) {
      // Return the first error if there are any.
      return { data, error: errors[0] };
    }

    return { data };
  } catch (error) {
    return { error };
  }
}

