<?php

namespace App\Modules\Image;

use Illuminate\Support\Facades\Storage;
use Intervention\Image\Facades\Image;

class Uploader
{
    /**
     * 图片路径
     *
     * @var string
     */
    public $path;

    /**
     * 图片名称
     *
     * @var string
     */
    public $filename;

    /**
     * 保存路径
     *
     * @var string
     */
    public $savePath;

    /**
     * 图片实例
     *
     * @var \Intervention\Image\Image
     */
    protected $img;

    /**
     * 裁剪X坐标点
     *
     * @var int
     */
    protected $croppedX;

    /**
     * 裁剪Y坐标点
     *
     * @var int
     */
    protected $croppedY;

    /**
     * 裁剪宽度
     *
     * @var int
     */
    protected $croppedWidth;

    /**
     * 裁剪高度
     *
     * @var int
     */
    protected $croppedHeight;

    /**
     * 调整宽度
     *
     * @var int
     */
    protected $resizeWidth;

    /**
     * 调整高度
     *
     * @var int
     */
    protected $resizeHeight;

    /**
     * construct function
     *
     * @param string $path 图片路径
     * @param string $filename 文件名称
     * @param string $savePath 保存路径
     */
    public function __construct($path, $filename, $savePath = 'images')
    {
        $this->path = $path;
        $this->filename = $filename;
        $this->savePath = $savePath;
        $this->img = Image::make($this->path);
    }

    /**
     * 检查make图片的结果是否成功，如果失败，后续的操作都无法进行
     *
     * @return boolean
     */
    public function imgMakeResult()
    {
        if (!$this->img) {
            return false;
        }
        return true;
    }

    /**
     * 保存图片
     *
     * @return boolean|string 成功返回图片链接，失败返回false
     */
    public function save()
    {
        if ($this->croppedWidth and $this->croppedHeight) {
            // 图片需要裁剪
            $this->img->crop(
                $this->croppedWidth,
                $this->croppedHeight,
                $this->croppedX,
                $this->croppedY
            );
        }

        if ($this->resizeWidth and $this->resizeHeight) {
            // 图片需要调整大小
            $this->img->resize($this->resizeWidth, $this->resizeHeight);
        }

        if (!Storage::disk('public')->exists($this->savePath)) {
            Storage::disk('public')->makeDirectory($this->savePath);
        }

        $relativePath = rtrim( $this->savePath, '/') . '/' . ltrim($this->filename, '/');
        $absolutePath = Storage::disk('public')->path($relativePath);
        $this->img->save($absolutePath);

        return Storage::url(ltrim($relativePath, '/'));
    }

    /**
     * 裁剪设置
     *
     * @param int $width 裁剪宽度
     * @param int $height 裁剪高度
     * @param int $x 裁剪X坐标点
     * @param int $y 裁剪Y坐标点
     * @return boolean
     */
    public function croppedSet($width, $height, $x, $y)
    {
        if (($x + $width) > $this->img->width() or ($y + $height) > $this->img->height()) {
            $this->img->destroy();
            return false;
        }

        $this->croppedX = $x;
        $this->croppedY = $y;
        $this->croppedWidth = $width;
        $this->croppedHeight = $height;
        return true;
    }

    /**
     * 宽高调整设置
     *
     * @param int $width 调整宽度
     * @param int $height 调整高度
     * @return void
     */
    public function resizeSet($width, $height)
    {
        $this->resizeWidth = $width;
        $this->resizeHeight = $height;
    }

    /**
     * 获取图片宽度
     *
     * @return int
     */
    public function width()
    {
        return $this->img->width();
    }

    /**
     * 获取图片高度
     *
     * @return int
     */
    public function height()
    {
        return $this->img->height();
    }
}
