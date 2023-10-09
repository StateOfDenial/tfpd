# Terraform Provider Docs
My first attempt at a Go CLI tool.

I get sick of having to go to https://registry.terraform.io in the browser 
so I made this CLI tool to go get the source in Markdown instead.

This tool will look for a `.terraform.lock.hcl` file and try to automatically 
decide what provider to look for in the docs. If more than one are present then
the list will be presented to you to pick from in a fuzzy finding TUI.

Then you will be presented with a list of options for the version to pick from
using the same fuzzy finder.

Finally, a list of possible documentation pages to look for (data, resource and 
others) will be presented in a fuzzy finder.

A lot of this is for practice and learning which is why many higher level 
libraries have not been used. Instead I opted to implement my own fuzzy finding
algorithm and a TUI using [tcell](https://github.com/gdamore/tcell). I may 
tcell with my own implementation later for fun but it's a lot of effort.

Also since I'm still learning, pretty sure the code will be very painful to
veterans.

## TODO
- [ ] Auto get provider version from `.terraform.lock.hcl`
- [ ] Prompt for provider to search for if none found
- [ ] Provide input arguments to shortcut steps
- [ ] Make TUI prettier 
- [ ] Implement some testing

## Contributions
Not going to look at contributions at this point, but feel free to fork it!
