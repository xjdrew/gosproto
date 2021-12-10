package sproto

/*
使用 go 语言类型表达 sproto 类型时，有些 sproto 的元信息（比如字段的tag值）较难直接表达。实现时，存在多种方案解决此问题。

考虑了很多因素后，gosproto 选择使用结构体 tag 能力，以类似注解(annotation)的形式描述 sproto 元信息。

meta tag 形式:
	sproto:"wire,tag,options..."
可选的 options 有 array, key, value 等。

wire: 字段的 sproto 类型
tag: 字段的 sproto tag 值
array: 存在此值，表示字段为数组
key=n: map 键的 sproto tag值
	若存在此值，
	编码时，将内存中的map转成数组，忽略key，直接将 map 的所有元素编码成数组
	解码时，使用数组生成map，map的键为数组元素中tag等于n的字段，map的值为数组元素本身
value=m：map 值的 sproto tag值
	若存在此值，则表示编解码时，将内存中的 map 结构处理成数组。
	数组元素类型为由 subtype 指定
	编码时，使用 map 的键值填充数组元素，按照tag值，分别填充对应字段
	解码时，使用数组的元素生成 map，map的键为数组元素tag等于n的字段，map的值为数组元素tag等于m的字段
subtype=field:
	当 value 值有效时，此值用来表示 map 转换为数组时的元素类型
*/
