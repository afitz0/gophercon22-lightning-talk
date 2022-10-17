# Temporal Demo
A simple demo app for my GopherCon 2022 lightning talk.

To run this demo like how it was run in that talk:

1. In one terminal, start a Temporal server. I used `temporalite` in ephemeral mode:
```
temporalite start --namespace default --ephemeral
```
  * For the demo, this terminal is never seen or used again.

1. In a different terminal, start a single worker:
```
go run worker/main.go
```

1. In a third terminal, start a workflow:
```
go run starter/main.go
```
  * In the first run, let it run to completion, watching or pointing out the output in the worker's terminal.

1. Open a browser to the Temporal UI. By default with `temporalite`, this is at http://localhost:8233/
  * By running with `--ephemeral`, there should be a single workflow.

1. Run the starter again, switching immediately back to the worker terminal. Once the second step ("Fulfilling order") starts, Ctrl-C kill it.

1. In yet another terminal (for those counting, this is the fourth), edit `workflow.go`.
  1. Mention: basically the same as what was in the slides.
  1. Make a token change or two to the code. I changed the final "Done!" println, as well as adding a println in the InitOrder func.
    * Goal: show that when the worker starts back up again that the new code is being run (not some cached something or other).

1. Open the UI again. There should be two workflows: the completed first one and one "Running...".
  * Clicking into the running workflow should show that it's stuck on the Fulfilling Order activity.

1. Start the worker back up
  * Note that the println added to InitOrder doesn't show up (nor does that activity's original "it's going to take N seconds")
  * Wait until completion. Note that the edited "Done!" message shows up, demonstrating that not only did we pick up with the last uncompleted activity, but it was running the new code.

1. If time permits, go back to the UI just to show that the server thinks everything happened successfully too.

The result would look something like this:

![Demo](slides/lightning-demo.gif)
