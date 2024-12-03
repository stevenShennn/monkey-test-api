# Input 板块需求

## 功能描述
Input 模块负责接收用户输入的 API 请求数据。该模块主要支持接收并解析用户提供的 cURL 格式的请求。后续版本可能会扩展支持其他输入格式（例如：原始 HTTP 请求、GraphQL 请求等）。

## 主要功能
1. **接收 cURL 格式输入**  
   用户可以将 cURL 请求直接粘贴到工具中，工具会自动解析并准备好相关参数供后续任务模块处理。

2. **支持批量输入**  
   允许用户一次性输入多个 cURL 请求，每个请求都将被独立拆解和处理。可以通过换行或特殊分隔符将多个 cURL 请求区分开。

3. **格式校验与解析**  
   输入的 cURL 请求需要经过严格的格式校验，确保其符合标准的 cURL 请求格式，并能正确解析 HTTP 方法、请求 URL、请求头、请求体等信息。

4. **扩展性**  
   当前仅支持 cURL 格式，后期可以根据需求扩展支持其他格式的 API 请求输入，例如：
    - 原始 HTTP 请求（支持 GET、POST、PUT、DELETE 等常用 HTTP 方法）。
    - 支持 JSON 格式输入，允许用户直接传递 JSON 格式的请求数据。
    - 支持 GraphQL 查询请求格式。

5. **支持代理设置**  
   在输入模块中，用户可以配置代理信息，允许通过代理服务器发送请求进行 API 测试。

6. **自动化处理**  
   对于每个输入的 cURL 请求，Input 模块将自动提取所有必要的信息（例如请求方法、URL、请求头、请求体、参数等），并将其传递给任务模块进行进一步处理。

## 输入样例

### cURL 示例
```bash
curl 'https://api.bilibili.com/x/kv-frontend/namespace/data?appKey=333.1333&nscode=0&versionId=1720688647891' \
  -H 'accept: */*' \
  -H 'accept-language: zh-CN,zh;q=0.9' \
  -H 'origin: https://www.bilibili.com' \
  -H 'priority: u=1, i' \
  -H 'referer: https://www.bilibili.com/video/BV1e3t4etExj/?spm_id_from=333.788.videopod.sections&vd_source=8b1d334ebaae98b258844334ab4126be' \
  -H 'sec-ch-ua: "Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "macOS"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-site' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'
```
### 多条请求示例
```
curl 'https://api.bilibili.com/x/player/wbi/playurl?avid=113148447169565&bvid=BV1e3t4etExj&cid=25661412680&qn=80&fnver=0&fnval=4048&fourk=1&gaia_source=&from_client=BROWSER&is_main_page=true&need_fragment=false&isGaiaAvoided=false&session=1d52e8c25f42e1dac6c68be0a1a8a1ea&voice_balance=1&web_location=1315873&dm_img_list=\[\]&dm_img_str=V2ViR0wgMS4wIChPcGVuR0wgRVMgMi4wIENocm9taXVtKQ&dm_cover_img_str=QU5HTEUgKEludGVsLCBBTkdMRSBNZXRhbCBSZW5kZXJlcjogSW50ZWwoUikgVUhEIEdyYXBoaWNzIDYzMCwgVW5zcGVjaWZpZWQgVmVyc2lvbilHb29nbGUgSW5jLiAoSW50ZW&dm_img_inter=%7B%22ds%22:\[\],%22wh%22:\[5506,4982,46\],%22of%22:\[433,866,433\]%7D&w_rid=2e6deefb46b2bd8ab9527a45de2513bf&wts=1732875990' \
  -H 'accept: application/json, text/plain, */*' \
  -H 'accept-language: zh-CN,zh;q=0.9' \
  -H $'cookie: buvid3=16CBB5D7-8019-190C-3700-7A619FD55CBC01374infoc; b_nut=1732853101; bsource=search_baidu; _uuid=FE44A83C-4A3B-D42F-C431-103103631139F501390infoc; buvid4=190FDB38-02F4-FAB0-A19C-F279EEAF2EF001930-024112904-8eJfRLsJprJ6qrbkTG1q8A%3D%3D; buvid_fp=6e776fc8cb99affd741b413a95befa01; enable_web_push=DISABLE; home_feed_column=5; CURRENT_FNVAL=4048; rpdid=|(k||RlYYmR|0J\'u~Jk)lJuk|; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzMxMTIzNTMsImlhdCI6MTczMjg1MzA5MywicGx0IjotMX0.p70lJ9r_CCoObESdtxgMAtqJPtqrZPOFA-8jUEyde7c; bili_ticket_expires=1733112293; SESSDATA=c63af5e1%2C1748405175%2C2eb66%2Ab2CjAR3AGOLFuAgABE6i5WOKTiMP-DOSIH-xvyFQfUyLPcabITs6JvS-a5VLS7RkRzgNASVk9MTU1qTERnYzY5Z2kwSFBKZHRRek1qR25iU1pWbFhBMGdob0RydEI4Z25pczk3LVBLQjNXY0tEcEIxOWZXLW9TcmNfRnpDQ2JYc3dxSmdGaEpFLWZ3IIEC; bili_jct=e88546c6493eca5f7065608c4beabc37; DedeUserID=435305731; DedeUserID__ckMd5=b2f09c84f30a704f; sid=67hsfo6o; CURRENT_QUALITY=80; header_theme_version=CLOSE; browser_resolution=2035-1160; b_lsid=93589256_1937763FE56' \
  -H 'origin: https://www.bilibili.com' \
  -H 'priority: u=1, i' \
  -H 'referer: https://www.bilibili.com/video/BV1e3t4etExj/?spm_id_from=333.788.videopod.sections&vd_source=8b1d334ebaae98b258844334ab4126be' \
  -H 'sec-ch-ua: "Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "macOS"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-site' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'

curl -X GET https://example.com/api/v1/users \
  -H "Authorization: Bearer your-token"
```

## JSON 格式拆解模板
每个 cURL 请求被解析后，将按照以下模板转换为 JSON 格式，便于存入 MongoDB。

```
{
  "request_id": "12345",               // 唯一请求ID
  "method": "POST",                     // HTTP 请求方法 (POST, GET, PUT, DELETE)
  "url": "https://example.com/api/v1/user", // 请求的 URL
  "headers": {                          // 请求头
    "Content-Type": "application/json",
    "Authorization": "Bearer your-token"
  },
  "body": {                             // 请求体 (如果是 POST 或 PUT 请求)
    "name": "John",
    "age": 30,
    "email": "john@example.com"
  },
  "params": {                           // 请求参数 (如果是 GET 请求或者有查询字符串)
    "id": "123"
  },
  "proxy": {                            // 代理设置（可选）
    "enabled": false,
    "proxy_url": "http://proxy.example.com"
  }
}
```
查询需要对时间戳进行倒序查询，而且查询是非常频繁的

## 技术需求
- cURL 解析：能够解析 cURL 请求的基本元素：HTTP 方法、请求 URL、请求头和请求体。
- 支持多种格式：后期支持 JSON、GraphQL、原生 HTTP 请求等格式的输入解析。
- 代理支持：允许用户通过配置文件或命令行参数设置代理服务器。
- 错误处理：对无效或格式错误的输入进行提示，确保输入数据的正确性。
- 并发处理：能够支持大规模输入请求的并发解析。

## 用户故事
- 作为开发人员，我希望能够方便地将 cURL 请求复制并粘贴到工具中，这样可以快速开始 API 测试。
- 作为开发人员，我希望能够输入多个 cURL 请求，工具会自动拆解并生成多个测试用例，而无需手动配置。
- 作为开发人员，我希望能通过代理服务器来进行 API 测试，这样可以确保接口在不同网络环境下的表现。

## 编码风格
- 项目文件夹不能新增，可以新增文件，编码风格按照依赖注入方式进行，使用golang interface进行解耦
- 重要的方法都需要进行英文备注
- input板块的代码都应该在input文件夹下面，具体拆分工作应该在curl等具体文件夹下面
- 要写出具体的单元测试