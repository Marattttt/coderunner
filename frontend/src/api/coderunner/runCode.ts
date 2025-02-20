import { z } from "zod";
import { CODERUNNER_API_URL } from "../../constants";

/*
 type RunResult struct {
	ExitCode   int           `json:"exitCode"`
	Stdout     string        `json:"stdout"`
	Stderr     string        `json:"stderr"`
	TimeTook   time.Duration `json:"timeTook"`
	TimeTookMs int64         `json:"timeTookMs"`
	Extra      *any          `json:"extra,omitempty"`
}
*/

const responseSchema = z.object({
	exitCode: z.number().int(),
	stdout: z.string(),
	stderr: z.string(),
	timeTookMs: z.number().int(),
})

export interface CodeRunResponce {
    exitCode: number;
    stdout: string;
    stderr: string;
    timeTookMs: number;
}

export default async function runCode(code: string, lang: string): Promise<CodeRunResponce> {
	const url = CODERUNNER_API_URL + '/run/' + lang;

	const body = {
		code: code,
		language: lang,
	}

	console.debug(`Fetching ${url} with params ${JSON.stringify(body)}`)

	const resp = await fetch(url, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(body)
	})

	const json = await resp.json();

	if (json.error) {
		throw new Error(`received from responce: ${JSON.stringify(json.error)}`)
	}

	const parsed = responseSchema.safeParse(json)
	if (!parsed.success) {
		console.error({
			msg: 'could not parse coderunner api response',
			resp: json,
		})
		throw parsed.error
	}

	return parsed.data
}
