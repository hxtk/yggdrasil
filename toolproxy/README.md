# Tool Proxy

This is designed to provide remote access to a specific set of tools
in a non-interactive fashion and with the possibility for fine-grained
access control policies, while generating clear audit trails.

This is a "harm reduction" mechanism for instances where manual commands
must be run against a production workload for whatever reason—perhaps
there was a use case missed by your automation or some disaster recovery
is needed that cannot be automated or your team has not finished the
initial attempt at automating things.

Ideally, this software should rarely if ever be used, but it provides
several advantages over allowing administrators interactive sessions with
the production workload because it makes it trivial to see exactly whom
has executed what commands and when. This discourages malicious actors
and allows you to track otherwise-undocumented changes by well-meaning
system administrators.

This tool is heavily inspired by the "safe proxy" case study from
"Building Secure and Reliable Systems" (See https://sre.google/books).
