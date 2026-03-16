package backend.mapper;

import backend.dto.InternshipDTO;
import backend.entity.InternshipEntity;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface InternshipMapper {

    @Mapping(source = "company.companyName", target = "companyName")
    InternshipDTO toDTO(InternshipEntity internship);

}
