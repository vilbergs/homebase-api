import {
  json,
  serve,
  validateRequest,
  PathParams
} from "https://deno.land/x/sift@0.2.0/mod.ts";

interface Temperature {
  value: number;
  unit: "CELSIUS" | "FAHRENHEIT";
}

interface Zone {
  name: string;
  temperature: Temperature;
  humidity: number;
}

serve({
  "/zones": handleZones,
  "/zones/:id": updateZone,
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
      (body as { name: string }).name
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

async function updateZone(request: Request, params?: PathParams) {
  const { error, body } = await validateRequest(request, {
    PUT: { body: ["temperature"] }
  });

  if (error) {
    return json({ error }, { status: 500 });
  }

  

  const query = `
    mutation($id: ID!, $name:String, $temperature: TemperatureInput, $humidity: Float) {
      partialUpdateZone(
        id: $id,
        data: { name: $name, temperature: { create: $temperature }, humidity: $humidity}
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

 // Handle POST requests.
 if (request.method === "PUT" && params?.id) {
  const { data, error} = await queryFauna(
    query,
    { id: params.id, ...body },
  );
  if (error) {
    return json({ error }, { status: 500 });
  }

  return json(data, { status: 200 });
}

return json({ error: "couldn't update the zone" }, { status: 500 });
}

/** Get all zones available in the database. */
async function getAllZones() {
  const query = `
    query {
      allZones {
        data {
          _id
          name
          temperature {
            value
            unit
          }
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
  } = await queryFauna(query, {});
  if (error) {
    return { error };
  }

  return { zones };
}

/** Create a new zone in the database. */
async function createZone(
  name: string
): Promise<{
  name?: string;
  temperature?: Temperature;
  humidity?: number;
  error?: string;
}> {
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
  } = await queryFauna(query, { name });
  if (error) {
    return { error };
  }

  return createZone; // Zone
}

async function queryFauna(
  query: string,
  variables: { [key: string]: unknown }
): Promise<{
  data?: any;
  error?: any;
}> {

  console.log(variables)
  // Grab the secret from the environment.
  const token = Deno.env.get("FAUNA_SECRET");
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
        "X-Schema-Preview": "partial-update-mutation"

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
