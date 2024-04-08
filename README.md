# ktalk

ktalk is kubectl plugin designed to convert natural language to a kubectl command via ChatGPT. 

With an input of a natural language command, such as "give me a list of containers in kube-system namespace", ktalk translates this intent into a `kubectl` command such as `kubectl get pods -n kube-system -o jsonpath='{.items[*].spec.containers[*].name}'`

# Usage

kubectl ktalk <message>

Example:

`kubectl ktalk give me the number of containers in the clusters that are running as root`

Are you sure want to execute the following command? Press Enter to execute this:  `kubectl get pods --all-namespaces -o jsonpath='{..securityContext.privileged}' | grep -o "true" | wc -l`