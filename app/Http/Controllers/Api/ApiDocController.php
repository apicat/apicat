<?php

namespace App\Http\Controllers\Api;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Validation\ValidationException;
use App\Http\Controllers\Controller;
use App\Modules\EditorJsonToHtml\Parser;
use App\Repositories\ApiDoc\MockPathRepository;
use App\Repositories\Project\ApiCommonUrlRepository;
use App\Repositories\ApiDoc\ApiDocRepository;
use App\Repositories\Project\DocShareRepository;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\Project\TreeCacheRepository;
use App\Repositories\DownloadRepository;
use App\Repositories\Import\BaseRepository as ImportBaseRepository;
use App\Repositories\Import\ApiCatRepository as ImportApiCatRepository;
use App\Repositories\Import\MarkdownRepository as ImportMarkdownRepository;
use App\Repositories\Import\PostmanRepository as ImportPostmanRepository;
use App\Repositories\Export\BaseRepository as ExportBaseRepository;
use App\Repositories\Export\ApiCatRepository as ExportApiCatRepository;
use App\Repositories\Export\PdfRepository as ExportPdfRepository;
use App\Repositories\Export\PostmanRepository as ExportPostmanRepository;
use App\Jobs\MarkdownImport;
use App\Jobs\ApiCatJsonImport;
use App\Jobs\PostmanImport;
use App\Jobs\ApiCatJsonExport;
use App\Jobs\PdfExport;
use App\Jobs\PostmanExport;

class ApiDocController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api', 'in.this.project']);
    }

    /**
     * 创建文档
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function create(Request $request)
    {
        $request->validate([
            'parent_id' => ['nullable', 'integer', 'min:1'],
            'title' => ['nullable', 'string', 'max:255'],
            'iteration_id' => ['nullable', 'integer', 'min:1']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法创建文档',
            ]);
        }

        $parentID = $request->input('parent_id') ? $request->input('parent_id') : 0;
        $title = $request->input('title') ? $request->input('title') : '无标题';
        $content = '{"type":"doc","content":[{"type":"paragraph"}]}';
        $iterationId = $request->input('iteration_id') ?? 0;

        if (!$node = ApiDocRepository::addDoc(ProjectRepository::active()->id, $parentID, $title, $content, Auth::id(), $iterationId)) {
            throw ValidationException::withMessages([
                'result' => '文档创建失败，请稍后重试。',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '文档创建成功',
            'data' => [
                'id' => $node->id,
                'parent_id' => $node->parent_id,
                'title' => $node->title,
                'sub_nodes' => []
            ]
        ];
    }

    /**
     * 创建HTTP API文档
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function httpTemplate(Request $request)
    {
        $request->validate([
            'parent_id' => ['nullable', 'integer', 'min:1'],
            'title' => ['nullable', 'string', 'max:255'],
            'iteration_id' => ['nullable', 'integer', 'min:1']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法创建文档',
            ]);
        }

        $parentID = $request->input('parent_id') ? $request->input('parent_id') : 0;
        $title = $request->input('title') ? $request->input('title') : '无标题';
        $content = '{"type":"doc","content":[{"type":"http_api_url","attrs":{"url":"https://api.example.com","path":"","method":2,"bodyDataType":2}},{"type":"heading","attrs":{"level":3},"content":[{"type":"text","text":"请求参数"}]},{"type":"http_api_request_parameter","attrs":{"request_header":{"params":[],"title":"Header 请求参数"},"request_body":{"params":[],"title":"Body 请求参数"},"request_query":{"params":[],"title":"Query 请求参数"}}},{"type":"heading","attrs":{"level":3},"content":[{"type":"text","text":"请求参数示例"}]},{"type":"code_block","attrs":{"language":"json"}},{"type":"http_status_code","attrs":{"intro":"Response Status Code:","code":200,"codeDesc":"OK"}},{"type":"heading","attrs":{"level":3},"content":[{"type":"text","text":"返回参数"}]},{"type":"http_api_response_parameter","attrs":{"response_header":{"params":[],"title":"返回头部"},"response_body":{"params":[],"title":"返回参数"}}},{"type":"heading","attrs":{"level":3},"content":[{"type":"text","text":"返回参数示例"}]},{"type":"code_block","attrs":{"language":"json"}},{"type":"paragraph"}]}';
        $iterationId = $request->input('iteration_id') ?? 0;

        if (!$node = ApiDocRepository::addDoc(ProjectRepository::active()->id, $parentID, $title, $content, Auth::id(), $iterationId)) {
            throw ValidationException::withMessages([
                'result' => '文档创建失败，请稍后重试。',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '文档创建成功',
            'data' => [
                'id' => $node->id,
                'parent_id' => $node->parent_id,
                'title' => $node->title,
                'sub_nodes' => []
            ]
        ];
    }

    /**
     * 编辑文档
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function update(Request $request)
    {
        $request->validate([
            'doc_id' => 'required|integer|min:1',
            'title' => 'required|string|max:255',
            'content' => 'required|string',
            'notification_list' => 'nullable|array'
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法编辑文档',
            ]);
        }

        $node = ApiDocRepository::getNode($request->input('doc_id'));
        if (!$node or $node->project_id != ProjectRepository::active()->id) {
            throw ValidationException::withMessages([
                'doc_id' => '文档不存在',
            ]);
        }

        if ($node->title != $request->input('title')) {
            $node->title = $request->input('title');
            TreeCacheRepository::remove($node->project_id);
        }

        $node->content = $request->input('content');
        $node->save();

        if ($node->content) {
            $content = json_decode($node->content, true);

            if (isset($content['content']) and is_array($content['content'])) {
                $httpApiUrlFinded = false;

                foreach ($content['content'] as $v) {
                    if ($v['type'] == 'api_url' or $v['type'] == 'http_api_url') {
                        ApiCommonUrlRepository::add($node->project_id, $v['attrs']['url']);
                    }

                    if ($v['type'] == 'http_api_url') {
                        $httpApiUrlFinded = true;
                        MockPathRepository::updatePath(ProjectRepository::active()->id, $node->id, $v['attrs']['path'], $v['attrs']['method']);
                    }
                }

                if (!$httpApiUrlFinded) {
                    MockPathRepository::del(ProjectRepository::active()->id, $node->id);
                }
            }
        }

        return [
            'status' => 0,
            'msg' => '保存成功',
            'data' => [
                'doc_id' => $node->id,
                'title' => $node->title
            ]
        ];
    }

    /**
     * 文档重命名
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function rename(Request $request)
    {
        $request->validate([
            'doc_id' => ['required', 'integer', 'min:1'],
            'title' => ['required', 'string', 'max:255']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法编辑文档',
            ]);
        }

        $node = ApiDocRepository::getNode($request->input('doc_id'));
        if (!$node or $request->input('project_id') != $node->project_id) {
            throw ValidationException::withMessages([
                'doc_id' => '文档不存在',
            ]);
        }

        ApiDocRepository::rename($node, $request->input('title'));

        return ['status' => 0, 'msg' => '修改成功'];
    }

    /**
     * 复制文档
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function copy(Request $request)
    {
        $request->validate([
            'doc_id' => ['required', 'integer', 'min:1'],
            'iteration_id' => ['nullable', 'integer', 'min:1']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法编辑文档',
            ]);
        }

        $node = ApiDocRepository::getNode($request->input('doc_id'));
        if (!$node or $request->input('project_id') != $node->project_id) {
            throw ValidationException::withMessages([
                'doc_id' => '文档不存在',
            ]);
        }

        $iterationId = $request->input('iteration_id') ?? 0;
        if (!$newNode = ApiDocRepository::copyNode($node, $iterationId)) {
            throw ValidationException::withMessages([
                'node_id' => '复制失败，请稍后重试。',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '复制成功',
            'data' => [
                'id' => $newNode->id,
                'parent_id' => $newNode->parent_id,
                'title' => $newNode->title,
                'sub_nodes' => []
            ]
        ];
    }

    /**
     * 删除文档
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function remove(Request $request)
    {
        $request->validate([
            'doc_id' => ['required', 'integer', 'min:1'],
            'iteration_id' => ['nullable', 'integer', 'min:1']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法删除文档',
            ]);
        }

        $node = ApiDocRepository::getNode($request->input('doc_id'));
        if (!$node or $node->project_id != $request->input('project_id')) {
            throw ValidationException::withMessages([
                'doc_id' => '文档不存在',
            ]);
        }

        $iterationId = $request->input('iteration_id') ?? 0;
        if (!ApiDocRepository::removeDoc($node, $iterationId)) {
            throw ValidationException::withMessages([
                'result' => '删除失败，请稍后重试。',
            ]);
        }

        MockPathRepository::del(ProjectRepository::active()->id, $node->id);
        DocShareRepository::remove($node->project_id, Auth::id(), $node->id);

        return ['status' => 0, 'msg' => '删除成功'];
    }

    /**
     * 文档分享详情
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function shareDetail(Request $request)
    {
        $request->validate([
            'doc_id' => ['required', 'integer', 'min:1']
        ]);

        $node = ApiDocRepository::getNode($request->input('doc_id'));
        if (!$node or $request->input('project_id') != $node->project_id) {
            throw ValidationException::withMessages([
                'doc_id' => '文档不存在',
            ]);
        }

        if (ProjectRepository::active()->visibility) {
            // 公开项目
            return [
                'status' => 0,
                'msg' => '',
                'data' => [
                    'visibility' => 1,
                    'link' => route('app.detail', ['projectID' => ProjectRepository::active()->id, 'docID' => $node->id]),
                    'secret_key' => ''
                ]
            ];
        }

        $share = DocShareRepository::get(ProjectRepository::active()->id, Auth::id(), $node->id);
        if (!$share) {
            return [
                'status' => 0,
                'msg' => '',
                'data' => [
                    'visibility' => 0,
                    'link' => '',
                    'secret_key' => ''
                ]
            ];
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'visibility' => 0,
                'link' => route('doc.index', ['docID' => $share->doc_id]),
                'secret_key' => $share->secret_key
            ]
        ];
    }

    /**
     * 分享文档
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function share(Request $request)
    {
        $request->validate([
            'doc_id' => ['required', 'integer', 'min:1'],
            'share' => ['required', 'boolean']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法分享文档',
            ]);
        }

        $node = ApiDocRepository::getNode($request->input('doc_id'));
        if (!$node or $request->input('project_id') != $node->project_id) {
            throw ValidationException::withMessages([
                'doc_id' => '文档不存在',
            ]);
        }

        if ($request->input('share')) {
            if (!$share = DocShareRepository::create(ProjectRepository::active()->id, Auth::id(), $node->id)) {
                throw ValidationException::withMessages([
                    'result' => '分享失败，请稍后重试。',
                ]);
            }

            return [
                'status' => 0,
                'msg' => '开启分享成功',
                'data' => [
                    'link' => route('doc.index', ['docID' => $share->doc_id]),
                    'secret_key' => $share->secret_key
                ]
            ];
        } else {
            if (!DocShareRepository::remove(ProjectRepository::active()->id, Auth::id(), $node->id)) {
                throw ValidationException::withMessages([
                    'result' => '关闭分享失败，请稍后重试。',
                ]);
            }

            return [
                'status' => 0,
                'msg' => '关闭分享成功'
            ];
        }
    }

    /**
     * 重置分享文档的访问密码
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function shareSecretKey(Request $request)
    {
        $request->validate([
            'doc_id' => ['required', 'integer', 'min:1']
        ]);

        $secretKey = DocShareRepository::changeSharePassword($request->input('project_id'), Auth::id(), $request->input('doc_id'));
        if (!$secretKey) {
            throw ValidationException::withMessages([
                'result' => '修改失败，请稍后重试。',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '修改成功',
            'data' => $secretKey
        ];
    }

    public function import(Request $request)
    {
        $request->validate([
            'filename' => 'required|string|max:255',
            'file' => 'required|string|max:255',
            'type' => 'required|string|in:markdown,apicat,postman',
            'parent_id' => 'nullable|integer|min:0',
            'iteration_id' => 'nullable|integer|min:1',
        ]);

        $parentID = $request->input('parent_id') ? $request->input('parent_id') : 0;

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法导入文件到项目中',
            ]);
        }

        if ($request->input('type') == 'markdown') {
            $import = new ImportMarkdownRepository;

            $fileNamrArr = explode('.', $request->input('filename'));
            $originFileName = array_shift($fileNamrArr);

            $jobID = $import->initJob([
                'userID' => Auth::id(),
                'projectID' => $request->input('project_id'),
                'fileName' => $request->input('file'),
                'originFileName' => $originFileName,
                'parentID' => $parentID,
                'iterationID' => $request->input('iteration_id') ?? 0
            ]);

            if (!$jobID) {
                throw ValidationException::withMessages([
                    'filename' => '导入失败，请重试',
                ]);
            }

            MarkdownImport::dispatch($jobID);
        } elseif ($request->input('type') == 'postman') {
            $import = new ImportPostmanRepository;

            $jobID = $import->initJob([
                'userID' => Auth::id(),
                'projectID' => $request->input('project_id'),
                'fileName' => $request->input('file'),
                'parentID' => $parentID,
                'iterationID' => $request->input('iteration_id') ?? 0
            ]);

            if (!$jobID) {
                throw ValidationException::withMessages([
                    'filename' => '导入失败，请重试',
                ]);
            }

            PostmanImport::dispatch($jobID);
        } else {
            $import = new ImportApiCatRepository;

            $jobID = $import->initJob([
                'userID' => Auth::id(),
                'projectID' => $request->input('project_id'),
                'fileName' => $request->input('file'),
                'parentID' => $parentID,
                'iterationID' => $request->input('iteration_id') ?? 0
            ]);

            if (!$jobID) {
                throw ValidationException::withMessages([
                    'filename' => '导入失败，请重试',
                ]);
            }

            ApiCatJsonImport::dispatch($jobID);
        }

        return response()->json([
            'status' => 0,
            'msg' => '',
            'data' => $jobID
        ]);
    }

    public function importResult(Request $request)
    {
        $request->validate([
            'job_id' => 'required|string|max:255'
        ]);

        $import = new ImportBaseRepository;
        $result = $import->jobResult($request->input('job_id'));

        if ($result['status'] == ImportBaseRepository::STATUS_FAIL) {
            return response()->json([
                'status' => -1,
                'msg' => $result['msg'],
                'data' => [
                    'result' => $result['status'],
                    'description' => $result['msg']
                ]
            ]);
        }

        return response()->json([
            'status' => 0,
            'msg' => '',
            'data' => [
                'result' => $result['status'],
                'description' => $result['msg']
            ]
        ]);
    }

    public function export(Request $request)
    {
        $request->validate([
            'doc_id' => 'nullable|integer|min:1',
            'type' => 'required|string|in:pdf,apicat,postman'
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法导出文件',
            ]);
        }

        if ($request->input('type') == 'pdf') {
            $export = new ExportPdfRepository;

            $jobID = $export->initJob([
                'userID' => Auth::id(),
                'projectID' => $request->input('project_id'),
                'docID' => $request->input('doc_id')
            ]);

            if (!$jobID) {
                throw ValidationException::withMessages([
                    'filename' => '导出失败，请重试',
                ]);
            }

            PdfExport::dispatch($jobID);
        } elseif ($request->input('type') == 'postman') {
            $export = new ExportPostmanRepository;

            $jobID = $export->initJob([
                'userID' => Auth::id(),
                'projectID' => $request->input('project_id'),
                'docID' => $request->input('doc_id')
            ]);

            if (!$jobID) {
                throw ValidationException::withMessages([
                    'filename' => '导出失败，请重试',
                ]);
            }

            PostmanExport::dispatch($jobID);
        } else {
            $export = new ExportApiCatRepository;

            $jobID = $export->initJob([
                'userID' => Auth::id(),
                'projectID' => $request->input('project_id'),
                'docID' => $request->input('doc_id')
            ]);

            if (!$jobID) {
                throw ValidationException::withMessages([
                    'filename' => '导出失败，请重试',
                ]);
            }

            ApiCatJsonExport::dispatch($jobID);
        }

        return response()->json([
            'status' => 0,
            'msg' => '',
            'data' => $jobID
        ]);
    }

    public function exportResult(Request $request)
    {
        $request->validate([
            'job_id' => 'required|string|max:255'
        ]);

        $export = new ExportBaseRepository;
        $result = $export->jobResult($request->input('job_id'));

        if ($result['status'] == ExportBaseRepository::STATUS_FAIL) {
            return response()->json([
                'status' => -1,
                'msg' => $result['msg'],
                'data' => [
                    'result' => $result['status'],
                    'description' => $result['msg']
                ]
            ]);
        }

        if ($result['status'] == ExportBaseRepository::STATUS_WAIT) {
            return response()->json([
                'status' => 0,
                'msg' => $result['msg'],
                'data' => [
                    'result' => $result['status'],
                    'description' => $result['msg']
                ]
            ]);
        }

        $downloadRepository = new DownloadRepository();
        if (!$downloadUrl = $downloadRepository->url($result['file'], $result['fileType'], $result['exportFileName'])) {
            throw ValidationException::withMessages([
                'result' => '导出失败，请稍后重试。',
            ]);
        }

        return response()->json([
            'status' => 0,
            'msg' => '',
            'data' => [
                'result' => $result['status'],
                'description' => $result['msg'],
                'url' => $downloadUrl
            ]
        ]);
    }
}
