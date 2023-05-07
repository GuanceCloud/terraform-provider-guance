resource "guance_pipeline" "demo" {
  name     = "oac-demo"
  category = "logging"
  source   = [
    "nginx"
  ]
  is_default = false
  is_force   = false

  content = <<EOF
    add_pattern("date2", "%%{YEAR}[./]%%{MONTHNUM}[./]%%{MONTHDAY} %%{TIME}")
    
    # access log
    grok(_, "%%{NOTSPACE:client_ip} %%{NOTSPACE:http_ident} %%{NOTSPACE:http_auth} \\[%%{HTTPDATE:time}\\] \"%%{DATA:http_method} %%{GREEDYDATA:http_url} HTTP/%%{NUMBER:http_version}\" %%{INT:status_code} %%{INT:bytes}")
    
    # access log
    add_pattern("access_common", "%%{NOTSPACE:client_ip} %%{NOTSPACE:http_ident} %%{NOTSPACE:http_auth} \\[%%{HTTPDATE:time}\\] \"%%{DATA:http_method} %%{GREEDYDATA:http_url} HTTP/%%{NUMBER:http_version}\" %%{INT:status_code} %%{INT:bytes}")
    grok(_, '%%{access_common} "%%{NOTSPACE:referrer}" "%%{GREEDYDATA:agent}"')
    user_agent(agent)
    
    # error log
    grok(_, "%%{date2:time} \\[%%{LOGLEVEL:status}\\] %%{GREEDYDATA:msg}, client: %%{NOTSPACE:client_ip}, server: %%{NOTSPACE:server}, request: \"%%{DATA:http_method} %%{GREEDYDATA:http_url} HTTP/%%{NUMBER:http_version}\", (upstream: \"%%{GREEDYDATA:upstream}\", )?host: \"%%{NOTSPACE:ip_or_host}\"")
    grok(_, "%%{date2:time} \\[%%{LOGLEVEL:status}\\] %%{GREEDYDATA:msg}, client: %%{NOTSPACE:client_ip}, server: %%{NOTSPACE:server}, request: \"%%{GREEDYDATA:http_method} %%{GREEDYDATA:http_url} HTTP/%%{NUMBER:http_version}\", host: \"%%{NOTSPACE:ip_or_host}\"")
    grok(_,"%%{date2:time} \\[%%{LOGLEVEL:status}\\] %%{GREEDYDATA:msg}")
    
    group_in(status, ["warn", "notice"], "warning")
    group_in(status, ["error", "crit", "alert", "emerg"], "error")
    
    cast(status_code, "int")
    cast(bytes, "int")
    
    group_between(status_code, [200,299], "OK", status)
    group_between(status_code, [300,399], "notice", status)
    group_between(status_code, [400,499], "warning", status)
    group_between(status_code, [500,599], "error", status)
    
    
    nullif(http_ident, "-")
    nullif(http_auth, "-")
    nullif(upstream, "")
    default_time(time)
    EOF

  test_data = <<EOF
    127.0.0.1 - - [24/Mar/2021:13:54:19 +0800] "GET /basic_status HTTP/1.1" 200 97 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.72 Safari/537.36"
    EOF
}
