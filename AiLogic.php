<?php

namespace app\api\logic;

use app\common\basics\Logic;
use app\common\model\ai\AiChatLog;
use app\common\model\ai\AiToWork;
use app\common\server\ConfigServer;
use app\common\server\FileServer;
use app\common\server\JsonServer;
use app\common\server\storage\Driver;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\GuzzleException;
use GuzzleHttp\Psr7\Request;
use GuzzleHttp\Psr7\Utils;

class AiLogic extends Logic
{
    private const QW_KEY = 'sk-OTVwdAIbvI';
    private const QW_OPEN_URL = 'https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions';
    private const QW_IMAGE_URL = 'https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis';
    private const QW_VIDEO_URL = 'https://dashscope.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis';

    private $url = '';
    private $model = '';
    private $system = '';
    private $search = false;
    private $chat = '';
    private $reasoning_content = '';


    /**
     * @note 获取AI生图记录
     * @param  $parmas array 接收的参数
     * @param $uid string 用户的秘钥ID
     * @return array|\think\response\Json
     */
    public function getAiWorkList($parmas, $uid)
    {
        $type = $parmas['type'];
        $page = $parmas['page'];
        $pageSize = 10;
        $model = AiToWork::where('user_id', $uid)->where('status', 'in', [0, 1]);

        if ($type != 3) {
            $model->where('type', $type);
        }
        if (!empty($parmas['close_person_upload'])) {
            $model->whereRaw("TRIM(params) <> '{}'"); // 不展示个人上传的图片和视频
        }
        $model = $model->order('id', 'desc');
        //计算总页数
        $total = $model->count();

        $show = ($total === 0) ? 0 : ceil($total / $pageSize);
        $model = $model->page($page, $pageSize)->select();
        $list = [];
        foreach ($model as $item) {
            $list[] = [
                'task_id' => $item->task_id,
                'params' => json_decode($item->params),
                'type' => $item->type,
                'status' => $item->status,
                'work' => json_decode($item->work)
            ];
        }
        return JsonServer::success('获取成功', [
            'page' => $page,
            'page_count' => $show,
            'total' => $total,
            'list' => $list
        ]);
    }

    public function sendStream2($parmas, $uid)
    {
        // 允许脚本在客户端断开连接后继续执行
        ignore_user_abort(true);

        $model = $parmas['model'];
        $chat_id = $parmas['chat_id'];
        $msg = $parmas['msg'];
        $restart = $parmas['restart'];
        $is_deep_reflection = $parmas['is_deep_reflection'] ?? 0; // 是发深度思索，默认否
        $chat = $this->getChatLog($chat_id, $uid, $model);
        if ($chat === false) {
            return JsonServer::error('聊天记录不存在');
        }

        $this->initializeModel($model, $is_deep_reflection);

        $messages = $this->prepareMessages($chat, $msg, $restart);
        // 删除AI设置
        $chatMessages = array_slice($messages, 1);
        $chat->chat = json_encode($chatMessages);

        $json = $this->prepareJsonPayload($messages);
        $header = $this->prepareHeaders();
        $client = new Client();
        $request = new Request('POST', $this->url, $header, json_encode($json));
        $promise = $client->sendAsync($request, ['stream' => true]);
        $isStreamClosed = false;
        $isCheck = true;

        $promise->then(
            function ($response) use ($chat, &$isStreamClosed) {
                $body = $response->getBody();
                while (!$body->eof()) {
                    if (connection_aborted()) {
                        if (!$isStreamClosed) {
                            $this->streamClose($chat);
                            $isStreamClosed = true;
                        }
                        break;
                    }
                    $res = Utils::readLine($body);
                    $this->handleResponseData($res);
                }
            },
            function ($exception) use (&$isStreamClosed, $uid, $msg, &$isCheck) {
                if (!$isStreamClosed) {
                    $isCheck = false;
                    $this->streamError($exception->getResponse(), $uid, $msg);
                    $isStreamClosed = false;
                }
            }
        )->wait();

        if (!$isStreamClosed) {
            $this->streamClose($chat);
            $isStreamClosed = true;
        }

        // 注释：ES搜索功能 - 与AI核心功能无关，可选功能
        // if ($isCheck) {
        //     if ($parmas['model'] == 1) {
        //         $esService = new ElasticsearchService();
        //         // 插入关键词到 ES
        //         $esService->addKeyword($parmas['msg']);
        //     }
        // }
        return 'data: [LOG_ID]:' . $chat->id . "\n\n" . 'data: [DONE]';
    }

    private function getChatLog($chat_id, $uid, $model)
    {
        if ($chat_id) {
            $chat = AiChatLog::where('id', $chat_id)->where('user_id', $uid)->where('model_id', $model)->find();
            if (!$chat) {
                return false;
            }
        } else {
            $chat = new AiChatLog();
            $chat->user_id = $uid;
            $chat->model_id = $model;
            $chat->chat = '[]';
            $chat->save();
        }
        return $chat;
    }

    private function initializeModel($model, $is_deep_reflection = 0)
    {
        switch ($model) {
            case 0:
                $this->url = self::QW_OPEN_URL;
                $this->model = 'deepseek-v3';
                $this->system = '你现在是AiChat AI咨询助手，来自AiChat私域聊天平台，是AiChat AI中的一种，请以AiChat AI咨询助手的身份进行回答';
                break;
            case 1:
                $this->url = self::QW_OPEN_URL;
                $this->model = $is_deep_reflection ? 'qwq-plus' : 'qwen-max';
                $this->system = '你现在是AiChat AI搜索，来自AiChat私域聊天平台，是AiChat AI中的一种，请以AiChat AI搜索的身份在互联网上搜索相关答案并返回相关内容，如果用户输入了链接那么请拒绝用户的请求并告知用户我无法直接访问网页';
                $this->search = true;
                break;
            case 2:
                $this->url = self::QW_OPEN_URL;
                $this->model = 'deepseek-v3';
                $this->system = '你现在是AiChat AI家庭医生助手，来自AiChat私域聊天平台，是AiChat AI中的一种，请以AiChat AI家庭医生的身份进行各种医学知识相关的回答，以中医为主西医为辅';
                break;
            case 3:
                $this->url = self::QW_OPEN_URL;
                $this->model = 'deepseek-v3';
                $this->system = '你现在是AiChat AI宠物助手，来自AiChat私域聊天平台，是AiChat AI中的一种，请以AiChat AI宠物助手的身份进行各种宠物知识相关回答';
                break;
        }
        $this->system .= '，全程使用简体中文回答，如果回答中有数学相关公式请使用双$符加换行的markdown语法';
    }

    private function prepareMessages($chat, $msg, $restart)
    {
        $messages = [
            [
                'role' => 'system',
                'content' => $this->system
            ]
        ];
        if ($chat->exists()) {
            $chatMsg = json_decode($chat->chat, true);
            if (is_array($chatMsg)) {
                foreach ($chatMsg as $_msg) {
                    $messages[] = $_msg;
                }
            }
        }
        if ($restart == 1) {
            // 重新生成上一句ai的回答，那么就删除最后一句
            array_pop($messages);
        } else {
            $messages[] = [
                'role' => 'user',
                'content' => $msg
            ];
        }
        return $messages;
    }

    private function prepareJsonPayload($messages)
    {
        return [
            'model' => $this->model,
            'messages' => $messages,
            'stream' => true,
            'enable_search' => $this->search
        ];
    }

    private function prepareHeaders()
    {
        return [
            'Authorization' => 'Bearer ' . self::QW_KEY,
        ];
    }

    private function streamClose($chat)
    {
        // 处理流关闭后的逻辑，包括正常结束或者手动终止
        $msg = json_decode($chat->chat, true);
        $msg[] = [
            'role' => 'assistant',
            'reasoning_content' => $this->reasoning_content, // 深度思考的内容
            'content' => $this->chat // 只记录回答的结果
        ];
        $chat->chat = json_encode($msg);
        $chat->save();
        ob_flush();
        flush();
    }

    private function streamError($res, $uid, $msg)
    {
        // 处理流异常后的逻辑
        $error = json_decode($res->getBody()->getContents(), true);
        $errors = [
            'user_id' => $uid,
            'msg' => $msg,
            'error' => $error,
        ];
        $msg = $this->getAliyunErrorMsg($error['error']['code']);
        echo 'data:' . $msg . "\n\n";
        // 追加写入错误日志
        file_put_contents('error.txt', json_encode($errors, JSON_UNESCAPED_UNICODE), FILE_APPEND);
        ob_flush();
        flush();
    }

    private function handleResponseData($data)
    {
        $complete = json_decode($data);
        if (isset($complete->error)) {
            echo 'data:' . $complete->error->code . "\n\n";
            ob_flush();
            flush();
            exit;
        }
        $res = $this->parseData($data);
        $word = is_array($res) ? $res['data'] : $res;

        if (!empty($word)) {
            if ($word == 'data: [DONE]' || $word == 'data: [CONTINUE]') {
                ob_flush();
                flush();
            } else {
                // 将对话保存到变量中
                if (is_array($res)) {
                    $this->reasoning_content .= $word;
                } else {
                    $this->chat .= $word;
                }

                $word = str_replace("\n", '<br/>', $word);
                $word = str_replace(" ", '&nbsp;', $word);

                if (is_array($res)) {
                    $text = json_encode(['type' => 'reasoning_content', 'data' => $word], JSON_UNESCAPED_UNICODE);
                    echo "data:$text\n\n";
                } else {
                    echo "data:$word\n\n";
                }

                ob_flush();
                flush();
            }
        }
    }

    private function parseData($data)
    {
        $data = str_replace('data: {', '{', $data);
        $data = str_replace('data:{', '{', $data);
        $data = rtrim($data, "\n\n");
        if (strpos($data, "}\n\n{") !== false) {
            $arr = explode("}\n\n{", $data);
            $data = "{{$arr[1]}";
        }
        if (strpos($data, 'data: [DONE]') !== false) {
            return 'data: [DONE]';
        } else {
            $data = json_decode($data, true);
            if (!is_array($data)) {
                return '';
            }
            if (isset($data['error'])) {
                echo 'data:' . $this->getAliyunErrorMsg($data['error']['code']) . "\n\n";
                echo 'data: [EXCEPTION]';
                ob_flush();
                flush();
                exit;
            }
            switch ($data['choices']['0']['finish_reason']) {
                case 'stop':
                    return 'data: [DONE]';
                case 'length':
                    return 'data: [CONTINUE]';
            }

            $delta = $data['choices']['0']['delta'];

            // 优先显示深度思考过程
            if (!empty($delta['reasoning_content'])) {
                return ['data' => $delta['reasoning_content']];
            } else {
                return $delta['content'] ?? '';
            }
        }
    }

    /**
    * @note AI根据文字生图
    * @param  $params array接收的参数
    * @param $uid string 用户ID
    * @return array
    */
    public function toImage($parmas, $uid)
    {
        $model = $parmas['model'];
        $prompt = $parmas['prompt'];
        $size = $parmas['size'];
        $n = $parmas['n'];
        $watermark = $parmas['watermark'] == 1 ? true : false;

        $this->url = self::QW_IMAGE_URL;
        switch ($model) {
            case 1:
                $this->model = 'wanx2.1-t2i-turbo';
                break;
            case 2:
                $this->model = 'wanx2.1-t2i-plus';
                break;
        }

        $headers = [
            'X-DashScope-Async' => 'enable',
            'Content-Type' => 'application/json',
            'Authorization' => 'Bearer ' . self::QW_KEY
        ];

        $json = [
            'model' => $this->model,
            'input' => [
                'prompt' => $prompt,
            ],
            'parameters' => [
                'size' => "$size[0]*$size[1]",
                'n' => (int)$n,
                'watermark' => $watermark
            ]
        ];

        $client = new Client();
        try {
            $resp = $client->post($this->url, [
                'headers' => $headers,
                'json' => $json
            ]);
        } catch (GuzzleException $e) {
            $responseBody = $e->getResponse()->getBody()->getContents();
            $responseData = json_decode($responseBody, true);
            $errorCode = $responseData['code'] ?? 'Unknown error';
            $msg = $this->getAliyunErrorMsg($errorCode);
            return JsonServer::error("图片生成失败，$msg");
        }
        $body = $resp->getBody()->getContents();
        $data = json_decode($body, true);
        if (isset($data['output']['task_id'])) {
            // 图片生成中
            AiToWork::create([
                'user_id' => $uid ?? 0,
                'task_id' => $data['output']['task_id'],
                'params' => json_encode($parmas),
                'type' => 1,
                'status' => 0,
            ]);
            return JsonServer::success('图片生成中', $data['output']['task_id']);
        } else {
            return JsonServer::error('图片生成失败，' . $data['message']);
        }
    }

    /**
    * @note 获取配置的错误信息
    * @param $code int 状态码
    * @return string
    */
        
    private function getAliyunErrorMsg($code)
    {
        switch ($code) {
            case 'InvalidParameter':
                return 'AiChat AI服务器错误，请联系客服处理或稍后再试';
            case 'APIConnectionError':
                return 'AiChat AI服务器网络出错，请稍后再试';
            case 'invalid_request_error':
                return 'AiChat AI服务器错误，请联系客服处理或稍后再试';
            case 'DataInspectionFailed':
            case 'data_inspection_failed':
                return '敏感内容，生成失败';
            case 'IPInfringementSuspect':
                return 'IP侵权，生成失败';
            default:
                return '服务错误，请稍后再试';
        }
    }


    /**
    * @note 获取进度的任务列表
    * @param $task_id int 任务ID
    * @param $uid string 用户ID
    * @return string
    */
    private function getTask($task_id, $uid)
    {
        $model = AiToWork::where('task_id', $task_id)->where('user_id', $uid)->find();

        if (!$model) {
            return [
                'output' => [
                    'task_status' => 'NULL',
                    'message' => '任务不存在'
                ]
            ];
        }
        if ($model->status == 1) {
            return [
                'output' => [
                    'task_status' => 'OK',
                    'results' => json_decode($model->work, true)
                ]
            ];
        } else if ($model->status == 2) {
            return [
                'output' => [
                    'task_status' => 'FAILED',
                    'message' => json_decode($model->work, true)['error']
                ]
            ];
        }
        $url = "https://dashscope.aliyuncs.com/api/v1/tasks/$task_id";
        $headers = [
            'Authorization' => 'Bearer ' . self::QW_KEY
        ];
        $client = new Client();
        try {
            $resp = $client->get($url, [
                'headers' => $headers
            ]);
        } catch (GuzzleException $e) {
            $responseBody = $e->getResponse()->getBody()->getContents();
            $responseData = json_decode($responseBody, true);
            $errorCode = $responseData['code'] ?? 'Unknown error';
            $msg = $this->getAliyunErrorMsg($errorCode);
            return [
                'output' => [
                    'task_status' => 'FAILED',
                    'message' => $msg
                ]
            ];
        }
        $body = $resp->getBody()->getContents();
        $body = json_decode($body, true);
        if (isset($body['output']['task_status']) && $body['output']['task_status'] == 'UNKNOWN') {
            $model->work = json_encode([
                'error' => '任务不存在'
            ]);
            $model->status = 2;
            $model->save();
        }
        return $body;
    }


    /**
    * @note 获取图片ID
    * @param $params array 接收参数
    * @param $uid strsing 用户ID
    * @return array
    */
    public function getImage($parmas, $uid)
    {
        $task_id = $parmas['task_id'];
        $data = $this->getTask($task_id, $uid);
        if ($data['output']['task_status'] == 'OK') {
            return JsonServer::success('图片生成成功', [
                'iamge' => $data['output']['results'],
                'image' => $data['output']['results']
            ]);
        } else if ($data['output']['task_status'] == 'SUCCEEDED') {
            // 图片生成成功
            $result = $data['output']['results'];
            $url = [];
            $i = 1;
            foreach ($result as $item) {
                if (isset($item['url'])) {
                    $path = "uploads/ai_images/{$task_id}_$i.png";
                    FileServer::uploadUrl($item['url'], $path);
                } else {
                    $path = "images/aichat_uni/ai/ai_picture/icon_error.png";
                }
                $url[] = "https://static.jsss999.com/$path";
                $i++;
            }
            // 上传数据库
            $model = AiToWork::where('task_id', $task_id)->find();
            $model->work = json_encode($url);
            $model->status = 1;
            $model->save();

            return JsonServer::success('图片生成成功', [
                'iamge' => $url,
                'image' => $url
            ]);
        } else if ($data['output']['task_status'] == 'FAILED') {
            if (isset($data['output']['code'])) {
                $data['output']['message'] = $this->getAliyunErrorMsg($data['output']['code']);
            }
            // 图片生成失败
            AiToWork::where('task_id', $task_id)->update([
                'work' => json_encode([
                    'error' => $data['output']['message']
                ]),
                'status' => 2
            ]);
            return JsonServer::error('图片生成失败，' . $data['output']['message']);
        } else if ($data['output']['task_status'] == 'RUNNING' || $data['output']['task_status'] == 'PENDING' || $data['output']['task_status'] == 'SUSPENDED') {
            // 图片生成中
            return JsonServer::success('图片生成中', $task_id, 1, 1);
        } else {
            return JsonServer::error('任务不存在');
        }
    }

    /**
    * @note AI生视频
    * @param  $params array 接收的参数
    * @param $uid string 秘钥ID
    * @return array
    */
    public function toVideo($parmas, $uid)
    {
        $to = $parmas['to'];
        $prompt = $parmas['prompt'] ?? '';
        $img_url = $parmas['img_url'] ?? '';

        switch ($to) {
            case 1:
                $this->model = 'wanx2.1-i2v-plus';
                break;
            case 2:
                $this->model = 'wanx2.1-t2v-turbo';
                break;
        }
        $this->url = self::QW_VIDEO_URL;
        $headers = [
            'X-DashScope-Async' => 'enable',
            'Content-Type' => 'application/json',
            'Authorization' => 'Bearer ' . self::QW_KEY
        ];
        $input = [];
        if ($prompt) {
            $input['prompt'] = $prompt;
        }
        if ($img_url) {
            // 检查是否为本地/局域网URL，如果是则转换为base64
            if ($this->isLocalOrPrivateUrl($img_url)) {
                $input['img_url'] = $this->convertImageUrlToBase64($img_url);
            } else {
                $input['img_url'] = $img_url;
            }
        }
        $json = [
            'model' => $this->model,
            'input' => $input
        ];
        $client = new Client();
        try {
            $resp = $client->post($this->url, [
                'headers' => $headers,
                'json' => $json
            ]);
        } catch (GuzzleException $e) {
            $responseBody = $e->getResponse()->getBody()->getContents();
            $responseData = json_decode($responseBody, true);
            $errorCode = $responseData['code'] ?? 'Unknown error';
            $msg = $this->getAliyunErrorMsg($errorCode);
            return JsonServer::error("视频生成失败，$msg");
        }
        $body = $resp->getBody()->getContents();
        $data = json_decode($body, true);
        if (isset($data['output']['task_id'])) {
            // 视频生成中
            AiToWork::create([
                'user_id' => $uid ?? 0,
                'task_id' => $data['output']['task_id'],
                'params' => json_encode([
                    'to' => $to,
                    'prompt' => $prompt,
                    'img_url' => $img_url
                ]),
                'type' => 2,
                'status' => 0,
            ]);
            return JsonServer::success('视频生成中', $data['output']['task_id']);
        } else {
            return JsonServer::error('视频生成失败，' . $data['message']);
        }
    }

    /**
    * @note 查看视频生成情况1111
    * @param $params array 接收的参数
    * @param $uid stsring 用户的秘钥ID
    * @return array
    */
    public function getVideo($parmas, $uid)
    {
        $task_id = $parmas['task_id'];
        $data = $this->getTask($task_id, $uid);
        if ($data['output']['task_status'] == 'OK') {
            return JsonServer::success('视频生成成功', ['video' => $data['output']['results']]);
        } else if ($data['output']['task_status'] == 'SUCCEEDED') {
            $model = AiToWork::where('task_id', $task_id)->find();

            if ($model->status == 3) {
                return JsonServer::error('水印正在生成中');
            }
            // 视频生成成功
            if (!empty($data['output']['video_url'])) {
                $result = $data['output']['video_url'];
            } else {
                $result = $data['output']['results']['video_url'] ?? '';
            }
            if (empty($result)) {
                return JsonServer::error('数据错误，请联系管理员');
            }
            $localInputVideo = public_path() . "uploads/temp/{$task_id}.mp4";
            // 下载视频到本地
            $this->downloadVideo($result, $localInputVideo);
            $path = "uploads/ai_videos/{$task_id}.mp4";
            $url = [];
            $url[] = "https://static-nine-world.oss-cn-shanghai.aliyuncs.com/$path";

            //视频添加水印
            $paths = "uploads/temp/{$task_id}_" . date('YmdHis') . ".mp4";
            $localOutputVideo = public_path() . $paths;
            // 水印文字内容
            $watermarkImage = public_path() . '/static/ai_chat/222.png';
            $model->status = 3;
            $model->save();
            // 构建 FFmpeg 命令，将水印图片叠加到视频右下角
            $command = "ffmpeg -i $localInputVideo -i $watermarkImage -filter_complex \"overlay=W-w-10:H-h-10\" $localOutputVideo";
            // 执行命令
            exec($command, $output, $returnCode);

            if ($returnCode != 0) {
                $model->status = 2;
                $model->save();
                return JsonServer::error('水印添加失败', [], 1, 1);
            }

            FileServer::uploadUrl(public_path() . $paths, $path);
            // 上传数据库
            $model->work = json_encode($url);
            $model->status = 1;
            $model->save();

            unlink($localOutputVideo);
            unlink($localInputVideo);
            return JsonServer::success('视频生成成功', ['video' => $url]);
        } else if ($data['output']['task_status'] == 'FAILED') {
            if (isset($data['output']['code'])) {
                $data['output']['message'] = $this->getAliyunErrorMsg($data['output']['code']);
            }
            // 视频生成失败
            AiToWork::where('task_id', $task_id)->update([
                'work' => json_encode([
                    'error' => $data['output']['message']
                ]),
                'status' => 2
            ]);
            return JsonServer::error('视频生成失败，' . $data['output']['message']);
        } else if ($data['output']['task_status'] == 'RUNNING' || $data['output']['task_status'] == 'PENDING' || $data['output']['task_status'] == 'SUSPENDED') {
            // 视频生成中
            return JsonServer::success('视频生成中', $task_id, 1, 1);
        } else {
            return JsonServer::error('任务不存在');
        }
    }

    /**
     * 下载视频到本地
     * @param string $url 视频链接
     * @param string $filename 本地保存的文件名
     */
    private function downloadVideo($url, $filename)
    {
        //创建目录出来，如果没有创建的话
        $savePath = dirname($filename);
        if (!is_dir($savePath)) {
            createFolders($savePath);
        }
        $content = file_get_contents($url);
        file_put_contents($filename, $content);
    }

   
     /**
    * @note 设置阿里同步的acl
    * @param $path string path路径
    * @param $acl string 目录权限
    * @return bool|json
    */
    public function setAcl($path, $acl = 'public-read')
    {
        $config = [
            'default' => ConfigServer::get('storage', 'default', 'local'),
            'engine' => ConfigServer::get('storage_engine') ?? ['local' => []]
        ];

        $StorageDriver = new Driver($config);
        if ($StorageDriver->setAcl($path, $acl)) {
            return JsonServer::success('设置成功');
        } else {
            return JsonServer::error('设置失败');
        }
    }

    /**
     * @notes 检查URL是否为本地或局域网地址
     * @param string $url 要检查的URL
     * @return bool 是否为本地/局域网地址
     */
    private function isLocalOrPrivateUrl($url)
    {
        $parsedUrl = parse_url($url);
        if (!$parsedUrl || !isset($parsedUrl['host'])) {
            return false;
        }
        
        $host = $parsedUrl['host'];
        
        // 检查是否为localhost相关
        if (in_array($host, ['localhost', '127.0.0.1', '::1'])) {
            return true;
        }
        
        // 检查是否为局域网IP段
        if (filter_var($host, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4)) {
            // 私有IP段：10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
            $ip = ip2long($host);
            $privateRanges = [
                ['10.0.0.0', '10.255.255.255'],
                ['172.16.0.0', '172.31.255.255'],
                ['192.168.0.0', '192.168.255.255']
            ];
            
            foreach ($privateRanges as $range) {
                if ($ip >= ip2long($range[0]) && $ip <= ip2long($range[1])) {
                    return true;
                }
            }
        }
        
        return false;
    }

    /**
     * @throws \Exception
     * @notes 将本地图片URL转换为base64格式供阿里云API使用
     * @param string $imageUrl 本地图片URL地址
     * @return string base64格式的图片数据
     */
    private function convertImageUrlToBase64($imageUrl)
    {
        try {
            // 使用GuzzleHttp客户端获取图片内容
            $client = new Client(['timeout' => 30]);
            $response = $client->get($imageUrl);
            
            // 获取图片二进制数据
            $imageContent = $response->getBody()->getContents();
            
            // 获取Content-Type来确定MIME类型
            $contentType = $response->getHeaderLine('Content-Type');
            if (empty($contentType)) {
                // 如果没有Content-Type，根据URL扩展名判断
                $extension = strtolower(pathinfo(parse_url($imageUrl, PHP_URL_PATH), PATHINFO_EXTENSION));
                switch ($extension) {
                    case 'jpg'://测试代码
                    case 'jpeg':
                        $contentType = 'image/jpeg';
                        break;
                    case 'png':
                        $contentType = 'image/png';
                        break;
                    case 'gif':
                        $contentType = 'image/gif';
                        break;
                    case 'webp':
                        $contentType = 'image/webp';
                        break;
                    default:
                        $contentType = 'image/jpeg'; // 默认为jpeg
                }
            }
            
            // 转换为base64格式
            $base64Data = base64_encode($imageContent);
            
            // 返回符合阿里云API要求的格式：data:{MIME_type};base64,{base64_data}
            return "data:{$contentType};base64,{$base64Data}";
            
        } catch (\Exception $e) {
            throw new Exception('图片转换失败: ' . $e->getMessage());
        }
    }
}